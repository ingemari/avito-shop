package services

import (
	"avito-shop/internal/models"
	"avito-shop/internal/repositories"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(username, password string) (*models.User, error)
}

type authService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Login(username, password string) (*models.User, error) {
	if username == "" || password == "" {
		return nil, errors.New("empty user or pass")
	}
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		newUser := &models.User{
			Username: username,
			Password: string(hashedPassword),
			Balance:  1000,
		}
		if err := s.userRepo.CreateUser(newUser); err != nil {
			return nil, errors.New("failed to create user")
		}
		return newUser, nil
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
