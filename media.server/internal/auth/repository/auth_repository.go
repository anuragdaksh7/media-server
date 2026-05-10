package repository

import (
	"fileserver/internal/auth/model"

	"gorm.io/gorm"
)

type AuthRepository struct {
	DB *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{
		DB: db,
	}
}

func (r *AuthRepository) CreateUser(
	user *model.User,
) error {
	return r.DB.Create(user).Error
}

func (r *AuthRepository) FindByEmail(
	email string,
) (*model.User, error) {

	var user model.User

	err := r.DB.
		Where("email = ?", email).
		First(&user).
		Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}
