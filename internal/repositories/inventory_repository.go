package repositories

import (
	"avito-shop/internal/models"

	"gorm.io/gorm"
)

type InventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) *InventoryRepository {
	return &InventoryRepository{db: db}
}

func (r *InventoryRepository) AddItem(userID uint, itemType string, quantity int) error {
	var inventory models.Inventory

	if err := r.db.Where("user_id = ? AND item_type = ?", userID, itemType).First(&inventory).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// newItem
			newItem := models.Inventory{UserID: userID, ItemType: itemType, Quantity: quantity}
			return r.db.Create(&newItem).Error
		}
		return err
	}

	inventory.Quantity += quantity
	return r.db.Save(&inventory).Error
}

func (r *InventoryRepository) GetUserInventory(userID uint) ([]models.Inventory, error) {
	var inventory []models.Inventory
	err := r.db.Where("user_id = ?", userID).Find(&inventory).Error
	return inventory, err
}
