package mocks

import (
	"avito-shop/internal/models"

	"github.com/stretchr/testify/mock"
)

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) CreateTransaction(transaction *models.Transaction) error {
	args := m.Called(transaction)
	return args.Error(0)
}

func (m *MockTransactionRepository) GetUserTransactions(userID uint) ([]models.Transaction, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Transaction), args.Error(1)
}
