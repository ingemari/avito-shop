package services_test

import (
	"testing"

	"avito-shop/internal/models"
	"avito-shop/internal/services"
	"avito-shop/internal/services/mocks"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Login_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	authService := services.NewAuthService(mockRepo)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("hashedpassword"), bcrypt.DefaultCost)
	expectedUser := &models.User{
		ID:       1,
		Username: "testuser",
		Password: string(hashedPassword),
		Balance:  1000,
	}

	mockRepo.On("GetUserByUsername", "testuser").Return(expectedUser, nil)

	user, err := authService.Login("testuser", "hashedpassword")

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_EmptyPass(t *testing.T) {
	authService := services.NewAuthService(nil)
	_, err := authService.Login("", "")

	assert.Error(t, err)
	assert.Equal(t, "empty user or pass", err.Error())
}

func TestAuthService_Login_IncorrectPass(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	authService := services.NewAuthService(mockRepo)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
	user := &models.User{
		Username: "testuser",
		Password: string(hashedPassword),
		Balance:  1000,
	}

	mockRepo.On("GetUserByUsername", "testuser").Return(user, nil)

	_, err := authService.Login("testuser", "123")

	assert.Error(t, err)
	assert.Equal(t, "invalid credentials", err.Error())

	mockRepo.AssertExpectations(t)
}
