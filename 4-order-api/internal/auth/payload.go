package auth

type GenerateSessionIdRequest struct {
	PhoneNumber    string `json:"phoneNumber" validate:"required"`
}

type GenerateSessionIdResponse struct {
	SessionId string `json:"sessionId"`
}

type ValidationCodeRequest struct {
	SessionId string `json:"sessionId"`
	Code int `json:"validationCode"`
}

type ValidationCodeResponse struct {
	Token string `json:"token"`
}
