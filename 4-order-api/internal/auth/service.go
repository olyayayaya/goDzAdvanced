package auth

import (
	"dz4/internal/user"
	"fmt"
)

type AuthService struct {
	UserRepository *user.UserRepository
}

func NewAuthService(userRepository *user.UserRepository) *AuthService {
	return &AuthService{
		UserRepository: userRepository,
	}
}

func (service *AuthService) FindByPhoneNumber(phoneNumber string) bool { 
	existedUser, _ := service.UserRepository.FindByPhoneNumber(phoneNumber)
		return existedUser == nil
}

func (service *AuthService) Register(phoneNumber, sessionId string) {
	user := &user.User{
		PhoneNumber: phoneNumber,
		SessionId: sessionId,
	}
	_, err := service.UserRepository.Create(user)
	if err != nil {
		fmt.Println(err)
	}
}

func (service *AuthService) FindBySessionId(sessionId string) bool {
	existedUser, _ := service.UserRepository.FindBySessionId(sessionId)
	return existedUser == nil
}
