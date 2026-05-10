package service

import (
	"fileserver/internal/auth/dto"
	"fileserver/internal/auth/model"
	"fileserver/internal/auth/repository"
	"fileserver/internal/auth/utils"
	apperrors "fileserver/pkg/errors"

	"errors"

	"gorm.io/gorm"
)

type AuthService struct {
	Repo *repository.AuthRepository
}

func NewAuthService(
	repo *repository.AuthRepository,
) *AuthService {
	return &AuthService{
		Repo: repo,
	}
}

func (s *AuthService) Register(
	req dto.RegisterRequest,
) (string, error) {

	_, err := s.Repo.FindByEmail(req.Email)

	if err == nil {
		return "", apperrors.ErrUserAlreadyExists
	}

	hash, err := utils.HashPassword(req.Password)

	if err != nil {
		return "", err
	}

	user := model.User{
		Email:    req.Email,
		Password: hash,
		Role:     model.RoleUser,
	}

	err = s.Repo.CreateUser(&user)

	if err != nil {
		return "", err
	}

	return utils.GenerateJWT(
		user.ID,
		string(user.Role),
	)
}

func (s *AuthService) Login(
	req dto.LoginRequest,
) (string, error) {

	user, err := s.Repo.FindByEmail(req.Email)

	if err != nil {
		return "", apperrors.ErrInvalidCredentials
	}

	err = utils.CheckPassword(
		req.Password,
		user.Password,
	)

	if err != nil {
		return "", apperrors.ErrInvalidCredentials
	}

	return utils.GenerateJWT(
		user.ID,
		string(user.Role),
	)
}

func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
