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
	// Создаём мок
	mockRepo := new(mocks.MockUserRepository)
	authService := services.NewAuthService(mockRepo)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("hashedpassword"), bcrypt.DefaultCost)
	// Ожидаемый пользователь
	expectedUser := &models.User{
		ID:       1,
		Username: "testuser",
		Password: string(hashedPassword),
		Balance:  1000,
	}

	// Настройка мока
	mockRepo.On("GetUserByUsername", "testuser").Return(expectedUser, nil)

	// Вызов метода
	user, err := authService.Login("testuser", "hashedpassword")

	// Проверки
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login(t *testing.T) {
	// Создаём мок
	mockRepo := new(mocks.MockUserRepository)
	authService := services.NewAuthService(mockRepo)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("hashedpassword"), bcrypt.DefaultCost)
	// Ожидаемый пользователь
	expectedUser := &models.User{
		ID:       1,
		Username: "testuser",
		Password: string(hashedPassword),
		Balance:  1000,
	}

	// Настройка мока
	mockRepo.On("GetUserByUsername", "testuser").Return(expectedUser, nil)

	// Вызов метода
	user, err := authService.Login("testuser", "hashedpassword")

	// Проверки
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockRepo.AssertExpectations(t)
}
