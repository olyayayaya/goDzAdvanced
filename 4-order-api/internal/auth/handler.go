package auth

import (
	"dz4/configs"
	"dz4/internal/user"
	"dz4/pkg/jwt"
	"dz4/pkg/req"
	"dz4/pkg/res"
	"math/rand"
	"net/http"
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
		sessionId := handler.Generate()

		existedUser := handler.AuthService.FindByPhoneNumber(body.PhoneNumber)
		if !existedUser {
			handler.AuthService.Register(body.PhoneNumber, sessionId) // если юзер не существует, регистрируем с новым айди
		} else {
			handler.Update(body.PhoneNumber, sessionId) // если сущестует, перезаписываем айди
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

		existedUser := handler.AuthService.FindBySessionId(body.SessionId)
		if !existedUser && body.ValidationCode != 3245 {
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

func (handler *AuthHandler) Generate() string {
	var letterRunes = []rune("qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890")
	b := make([]rune, 16)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func (handler *AuthHandler) Update(phoneNumber, sessionId string) error {
	_, err := handler.UserRepository.Update(&user.User{
		PhoneNumber: phoneNumber,
		SessionId: sessionId,
	})
	if err != nil {
		return err
	}
	return nil
}
