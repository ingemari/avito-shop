package mocks

import (
	"avito-shop/internal/models"

	"github.com/stretchr/testify/mock"
)

// MockUserRepository реализует интерфейс UserRepository для тестирования
type MockUserRepository struct {
	mock.Mock
}

// GetUserByUsername возвращает пользователя по имени
func (m *MockUserRepository) GetUserByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	return args.Get(0).(*models.User), args.Error(1)
}

// GetUserByID возвращает пользователя по ID
func (m *MockUserRepository) GetUserByID(id uint) (*models.User, error) {
	args := m.Called(id)
	return args.Get(0).(*models.User), args.Error(1)
}

// CreateUser создаёт нового пользователя
func (m *MockUserRepository) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

// UpdateUser обновляет данные пользователя
func (m *MockUserRepository) UpdateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

// GetUserBalance возвращает баланс пользователя
func (m *MockUserRepository) GetUserBalance(userID uint) (float64, error) {
	args := m.Called(userID)
	return args.Get(0).(float64), args.Error(1)
}
