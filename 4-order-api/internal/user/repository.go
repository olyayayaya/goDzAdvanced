package user

import (
	"dz4/pkg/db"

	"gorm.io/gorm/clause"
)


type UserRepository struct {
	Database *db.Db
}

func NewUserRepository(database *db.Db) *UserRepository {
	return &UserRepository{Database: database}
}

func (repo *UserRepository) Create(user *User) (*User, error) {
	result := repo.Database.DB.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (repo *UserRepository) FindByPhoneNumber(phoneNumber string) (*User, error) {
	var user User
	result := repo.Database.DB.First(&user, "phone_number = ?", phoneNumber)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}


func (repo *UserRepository) FindBySessionId(sessionId string) (*User, error) {
	var user User
	result := repo.Database.DB.First(&user, "session_id = ?", sessionId)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (repo *UserRepository) Update(user *User) (*User, error) {
	result := repo.Database.DB.Clauses(clause.Returning{}).Where("phone_number = ?", user.PhoneNumber).Updates(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}