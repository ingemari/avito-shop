package mocks

import (
	"avito-shop/internal/models"

	"github.com/stretchr/testify/mock"
)

type MockInventoryRepository struct {
	mock.Mock
}

func (m *MockInventoryRepository) AddItem(userID uint, itemType string, quantity int) error {
	args := m.Called(userID, itemType, quantity)
	return args.Error(0)
}

func (m *MockInventoryRepository) GetUserInventory(userID uint) ([]models.Inventory, error) {
	args := m.Called(userID)
	if args.Get(0) != nil {
		return args.Get(0).([]models.Inventory), args.Error(1)
	}
	return nil, args.Error(1)
}
