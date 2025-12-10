package auth

import (
	"dz4/internal/user"
)

type AuthService struct {
	UserRepository *user.UserRepository
}

func NewAuthService(userRepository *user.UserRepository) *AuthService {
	return &AuthService{
		UserRepository: userRepository,
	}
}

func (service *AuthService) FindByPhoneNumber(phoneNumber string) error {
	_, err := service.UserRepository.FindByPhoneNumber(phoneNumber)
	return err
}

func (service *AuthService) Register(phoneNumber string, sessionId string, code int) {
	user := &user.User{
		PhoneNumber: phoneNumber,
		SessionId:   sessionId,
		Code:        code,
	}
	_, err := service.UserRepository.Create(user)
	if err != nil {
		return
	}
}

func (service *AuthService) FindBySessionId(sessionId string) (int, error) {
    user, err := service.UserRepository.FindBySessionId(sessionId)
    if err != nil {
        return 0, err
    }
    return user.Code, nil
}

func (service *AuthService) Update(phoneNumber string, sessionId string, code int) {
	user := &user.User{
		PhoneNumber: phoneNumber,
		SessionId:   sessionId,
		Code:        code,
	}
	_, err := service.UserRepository.Update(user)
	if err != nil {
		return
	}
}
