package user

import (
	"dz4/internal/models"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	PhoneNumber string `gorm:"index"`
	SessionId   string
	Code        int
	Orders      []models.Order
}
