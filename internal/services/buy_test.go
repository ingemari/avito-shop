package services_test

import (
	"avito-shop/internal/models"
	"avito-shop/internal/services"
	"avito-shop/internal/services/mocks"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPurchaseItem_ItemNotFound(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockItemRepo := new(mocks.MockItemRepository)
	mockInvRepo := new(mocks.MockInventoryRepository)
	purchaseService := services.NewPurchaseService(mockUserRepo, mockItemRepo, mockInvRepo)

	mockItemRepo.On("GetItemByName", "item1").Return(nil, errors.New("not found"))

	balance, err := purchaseService.PurchaseItem(1, "item1")

	assert.Error(t, err)
	assert.Equal(t, "item not found", err.Error())
	assert.Equal(t, 0, balance)
	mockItemRepo.AssertExpectations(t)
}

func TestPurchaseItem_InsufficientFunds(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockItemRepo := new(mocks.MockItemRepository)
	mockInvRepo := new(mocks.MockInventoryRepository)
	purchaseService := services.NewPurchaseService(mockUserRepo, mockItemRepo, mockInvRepo)

	item := &models.Item{Name: "item1", Price: 100}
	user := &models.User{ID: 1, Balance: 50}

	mockItemRepo.On("GetItemByName", "item1").Return(item, nil)
	mockUserRepo.On("GetUserByID", uint(1)).Return(user, nil)

	balance, err := purchaseService.PurchaseItem(1, "item1")

	assert.Error(t, err)
	assert.Equal(t, "insufficient funds", err.Error())
	assert.Equal(t, 50, balance)
	mockUserRepo.AssertExpectations(t)
}

func TestPurchaseItem_FailedToUpdateBalance(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockItemRepo := new(mocks.MockItemRepository)
	mockInvRepo := new(mocks.MockInventoryRepository)
	purchaseService := services.NewPurchaseService(mockUserRepo, mockItemRepo, mockInvRepo)

	item := &models.Item{Name: "item1", Price: 100}
	user := &models.User{ID: 1, Balance: 200}

	mockItemRepo.On("GetItemByName", "item1").Return(item, nil)
	mockUserRepo.On("GetUserByID", uint(1)).Return(user, nil)
	mockUserRepo.On("UpdateUser", mock.Anything).Return(errors.New("update failed"))

	balance, err := purchaseService.PurchaseItem(1, "item1")

	assert.Error(t, err)
	assert.Equal(t, "failed to update balance", err.Error())
	assert.Equal(t, 100, balance)
	mockUserRepo.AssertExpectations(t)
}

func TestPurchaseItem_FailedToAddItem(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockItemRepo := new(mocks.MockItemRepository)
	mockInvRepo := new(mocks.MockInventoryRepository)
	purchaseService := services.NewPurchaseService(mockUserRepo, mockItemRepo, mockInvRepo)

	item := &models.Item{Name: "item1", Price: 100}
	user := &models.User{ID: 1, Balance: 200}

	mockItemRepo.On("GetItemByName", "item1").Return(item, nil)
	mockUserRepo.On("GetUserByID", uint(1)).Return(user, nil)
	mockUserRepo.On("UpdateUser", mock.Anything).Return(nil)
	mockInvRepo.On("AddItem", user.ID, item.Name, 1).Return(errors.New("inventory failed"))

	balance, err := purchaseService.PurchaseItem(1, "item1")

	assert.Error(t, err)
	assert.Equal(t, "failed to add item to inventory", err.Error())
	assert.Equal(t, 100, balance)
	mockInvRepo.AssertExpectations(t)
}

func TestPurchaseItem_Success(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockItemRepo := new(mocks.MockItemRepository)
	mockInvRepo := new(mocks.MockInventoryRepository)
	purchaseService := services.NewPurchaseService(mockUserRepo, mockItemRepo, mockInvRepo)

	item := &models.Item{Name: "item1", Price: 100}
	user := &models.User{ID: 1, Balance: 200}

	mockItemRepo.On("GetItemByName", "item1").Return(item, nil)
	mockUserRepo.On("GetUserByID", uint(1)).Return(user, nil)
	mockUserRepo.On("UpdateUser", mock.Anything).Return(nil)
	mockInvRepo.On("AddItem", user.ID, item.Name, 1).Return(nil)

	balance, err := purchaseService.PurchaseItem(1, "item1")

	assert.NoError(t, err)
	assert.Equal(t, 100, balance)

	mockUserRepo.AssertExpectations(t)
	mockItemRepo.AssertExpectations(t)
	mockInvRepo.AssertExpectations(t)
}
