package models

type Inventory struct {
	ID       uint `gorm:"primaryKey"`
	UserID   uint `gorm:"index"`
	ItemType string
	Quantity int
}
