package repositories

import (
	"avito-shop/internal/models"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreateTransaction(transaction *models.Transaction) error
	GetUserTransactions(userID uint) ([]models.Transaction, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) CreateTransaction(transaction *models.Transaction) error {
	return r.db.Create(transaction).Error
}

func (r *transactionRepository) GetUserTransactions(userID uint) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Where("from_user = ? OR to_user = ?", userID, userID).Find(&transactions).Error
	return transactions, err
}
