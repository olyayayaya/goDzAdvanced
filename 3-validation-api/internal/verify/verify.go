package verify

import (
	"crypto/rand"
	"dz3/configs"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/jordan-wright/email"
)

type VerifyHandlerDeps struct {
	*configs.Config
}

type VerifyHandler struct {
	*configs.Config
}

func NewVerifyHandler(router *http.ServeMux, deps VerifyHandlerDeps) {
	handler := &VerifyHandler{
		Config: deps.Config,
	}
	router.HandleFunc("POST /send", handler.Send())
	router.HandleFunc("GET /verify/{hash}", handler.Verify())
}

func (handler *VerifyHandler) Send() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqBody Request
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		validate := validator.New()
		err = validate.Struct(reqBody)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		hashBytes := make([]byte, 16)
		rand.Read(hashBytes)
		hash := hex.EncodeToString(hashBytes)

		data := map[string]string{
			"email": reqBody.Email,
			"hash":  hash,
		}

		file, _ := json.MarshalIndent(data, "", "  ")
		filename := fmt.Sprintf("%sverify.json", data["hash"])
		_ = os.WriteFile(filename, file, 0644)

		e := email.NewEmail()
		e.From = fmt.Sprintf("Email Verifier <%s>", handler.Config.Email)
		e.To = []string{reqBody.Email}
		e.Subject = "Please confirm your email"
		e.Text = []byte(fmt.Sprintf("To verify, open this:\nhttp://localhost:8081/verify/%s", hash))

		auth := smtp.PlainAuth("", handler.Config.Email, handler.Config.Password, handler.Config.Address)
		err = e.Send(handler.Config.Address+":587", auth)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Verification email sent"))
	}
}

func (handler *VerifyHandler) Verify() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := r.PathValue("hash")
		filename := fmt.Sprintf("%sverify.json", hash)

		file, err := os.ReadFile(filename)
		if err != nil {
			http.Error(w, "No verification data found", http.StatusNotFound)
			return
		}

		var data map[string]string
		json.Unmarshal(file, &data)

		if data["hash"] == hash {
			w.Write([]byte("true"))
		} else {
			w.Write([]byte("false"))
		}
		os.Remove(filename)

	}
}
