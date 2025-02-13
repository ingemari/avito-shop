package services

import (
	"avito-shop/internal/models"
	"avito-shop/internal/repositories"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *repositories.UserRepository
}

func NewAuthService(userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Login(username, password string) (*models.User, error) {
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		hashedPassword := HashPassword(password)
		newUser := &models.User{
			Username: username,
			Password: hashedPassword,
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

func HashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}
