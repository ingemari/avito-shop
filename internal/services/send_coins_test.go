package services_test

import (
	"avito-shop/internal/models"
	"avito-shop/internal/services"
	"avito-shop/internal/services/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTransferCoins_TransferToSelf(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockTransactionRepo := new(mocks.MockTransactionRepository)
	transactionService := services.NewTransactionService(mockUserRepo, mockTransactionRepo)

	user := &models.User{ID: 1, Username: "sender", Balance: 500}

	mockUserRepo.On("GetUserByID", uint(1)).Return(user, nil)
	mockUserRepo.On("GetUserByUsername", "sender").Return(user, nil)

	err := transactionService.TransferCoins(1, "sender", 100)

	assert.Error(t, err)
	assert.Equal(t, "cannot transfer to yourself", err.Error())
	mockUserRepo.AssertExpectations(t)
}

func TestTransferCoins_InsufficientFunds(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockTransactionRepo := new(mocks.MockTransactionRepository)
	transactionService := services.NewTransactionService(mockUserRepo, mockTransactionRepo)

	fromUser := &models.User{ID: 1, Username: "sender", Balance: 50}
	toUser := &models.User{ID: 2, Username: "recipient", Balance: 100}

	mockUserRepo.On("GetUserByID", uint(1)).Return(fromUser, nil)
	mockUserRepo.On("GetUserByUsername", "recipient").Return(toUser, nil)

	err := transactionService.TransferCoins(1, "recipient", 100)

	assert.Error(t, err)
	assert.Equal(t, "insufficient funds", err.Error())
	mockUserRepo.AssertExpectations(t)
}

func TestTransferCoins_Success(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockTransactionRepo := new(mocks.MockTransactionRepository)
	transactionService := services.NewTransactionService(mockUserRepo, mockTransactionRepo)

	fromUser := &models.User{ID: 1, Username: "sender", Balance: 500}
	toUser := &models.User{ID: 2, Username: "recipient", Balance: 100}

	mockUserRepo.On("GetUserByID", uint(1)).Return(fromUser, nil)
	mockUserRepo.On("GetUserByUsername", "recipient").Return(toUser, nil)

	mockUserRepo.On("UpdateUser", mock.Anything).Return(nil).Twice()
	mockTransactionRepo.On("CreateTransaction", mock.Anything).Return(nil)

	err := transactionService.TransferCoins(1, "recipient", 100)

	assert.NoError(t, err)
	assert.Equal(t, 400, fromUser.Balance)
	assert.Equal(t, 200, toUser.Balance)

	mockUserRepo.AssertExpectations(t)
	mockTransactionRepo.AssertExpectations(t)
}
