package repositories

import (
	"avito-shop/internal/models"

	"gorm.io/gorm"
)

type InventoryRepository interface {
	AddItem(userID uint, itemType string, quantity int) error
	GetUserInventory(userID uint) ([]models.Inventory, error)
}

type inventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) InventoryRepository {
	return &inventoryRepository{db: db}
}

func (r *inventoryRepository) AddItem(userID uint, itemType string, quantity int) error {
	var inventory models.Inventory

	if err := r.db.Where("user_id = ? AND item_type = ?", userID, itemType).First(&inventory).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			newItem := models.Inventory{UserID: userID, ItemType: itemType, Quantity: quantity}
			return r.db.Create(&newItem).Error
		}
		return err
	}

	inventory.Quantity += quantity
	return r.db.Save(&inventory).Error
}

func (r *inventoryRepository) GetUserInventory(userID uint) ([]models.Inventory, error) {
	var inventory []models.Inventory
	err := r.db.Where("user_id = ?", userID).Find(&inventory).Error
	return inventory, err
}
