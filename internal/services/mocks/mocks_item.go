package mocks

import (
	"avito-shop/internal/models"

	"github.com/stretchr/testify/mock"
)

type MockItemRepository struct {
	mock.Mock
}

func (m *MockItemRepository) GetItemByName(name string) (*models.Item, error) {
	args := m.Called(name)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Item), args.Error(1)
	}
	return nil, args.Error(1)
}
