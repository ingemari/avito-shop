package repositories

import (
	"avito-shop/internal/models"

	"gorm.io/gorm"
)

type ItemsRepository interface {
	GetItemByName(name string) (*models.Item, error)
}

type itemsRepository struct {
	db *gorm.DB
}

func NewItemsRepository(db *gorm.DB) ItemsRepository {
	return &itemsRepository{db: db}
}

func (r *itemsRepository) GetItemByName(name string) (*models.Item, error) {
	var item models.Item
	if err := r.db.Where("name = ?", name).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
