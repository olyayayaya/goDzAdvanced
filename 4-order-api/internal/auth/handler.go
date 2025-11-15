package auth

import (
	"dz4/configs"
	"dz4/pkg/jwt"
	"dz4/pkg/req"
	"dz4/pkg/res"
	"math/rand"
	"net/http"
	"strconv"
)

type AuthHandlerDeps struct {
	*configs.Config
	*AuthService
}

type AuthHandler struct {
	*configs.Config
	*AuthService
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		Config:      deps.Config,
		AuthService: deps.AuthService,
	}
	router.HandleFunc("POST /auth/createSessionId", handler.CreateSessionId())
	router.HandleFunc("POST /auth/checkValidationCode", handler.CheckValidationCode())
}

func (handler *AuthHandler) CreateSessionId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[GenerateSessionIdRequest](&w, r)
		if err != nil {
			return
		}
		sessionId := handler.GenerateSessionId()
		code := handler.GenerateCode()

		err = handler.AuthService.FindByPhoneNumber(body.PhoneNumber)
		if err != nil {
			handler.AuthService.Register(body.PhoneNumber, sessionId, code) // если юзер не существует, регистрируем с новым айди и кодом
		} else {
			handler.AuthService.Update(body.PhoneNumber, sessionId, code) // если сущестует, перезаписываем айди
		}

		data := GenerateSessionIdResponse{
			SessionId: sessionId,
		}
		res.Json(w, data, 200)
	}
}

func (handler *AuthHandler) CheckValidationCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[ValidationCodeRequest](&w, r)
		if err != nil {
			return
		}

		originalCode, err := handler.AuthService.FindBySessionId(body.SessionId)
		if err != nil || originalCode != body.Code {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		token, err := jwt.NewJWT(handler.Config.Auth.Secret).Create(body.SessionId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := ValidationCodeResponse{
			Token: token,
		}
		res.Json(w, data, 200)
	}

}

func (handler *AuthHandler) GenerateSessionId() string {
	var letterRunes = []rune("qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890")
	b := make([]rune, 16)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func (handler *AuthHandler) GenerateCode() int {
	var letterRunes = []rune("1234567890")
	b := make([]rune, 4)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	intCode, _ := strconv.Atoi(string(b))
	return intCode
}
