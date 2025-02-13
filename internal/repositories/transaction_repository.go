package repositories

import (
	"avito-shop/internal/models"

	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) CreateTransaction(transaction *models.Transaction) error {
	return r.db.Create(transaction).Error
}

func (r *TransactionRepository) GetUserTransactions(userID uint) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Where("from_user = ? OR to_user = ?", userID, userID).Find(&transactions).Error
	return transactions, err
}
