package services

import (
	"avito-shop/internal/models"
	"avito-shop/internal/repositories"
	"errors"
)

type TransactionService struct {
	userRepo        *repositories.UserRepository
	transactionRepo *repositories.TransactionRepository
}

func NewTransactionService(userRepo *repositories.UserRepository, transactionRepo *repositories.TransactionRepository) *TransactionService {
	return &TransactionService{
		userRepo:        userRepo,
		transactionRepo: transactionRepo,
	}
}

func (s *TransactionService) TransferCoins(fromUserID uint, toUsername string, amount int) error {
	fromUser, err := s.userRepo.GetUserByID(fromUserID)
	if err != nil {
		return errors.New("sender not found")
	}

	toUser, err := s.userRepo.GetUserByUsername(toUsername)
	if err != nil {
		return errors.New("recipient not found")
	}

	if fromUser.ID == toUser.ID {
		return errors.New("cannot transfer to yourself")
	}

	if fromUser.Balance < amount {
		return errors.New("insufficient funds")
	}

	fromUser.Balance -= amount
	toUser.Balance += amount

	if err := s.userRepo.UpdateUser(fromUser); err != nil {
		return err
	}
	if err := s.userRepo.UpdateUser(toUser); err != nil {
		return err
	}

	transaction := &models.Transaction{
		FromUser: fromUser.ID,
		ToUser:   toUser.ID,
		Amount:   amount,
	}

	return s.transactionRepo.CreateTransaction(transaction)
}
