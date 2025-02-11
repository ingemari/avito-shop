package models

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"uniqueIndex"`
	Password string
	Balance  int
}

type Transaction struct {
	ID       uint `gorm:"primaryKey"`
	FromUser uint `gorm:"index"`
	ToUser   uint `gorm:"index"`
	Amount   int
}

type Inventory struct {
	ID       uint `gorm:"primaryKey"`
	UserID   uint `gorm:"index"`
	ItemType string
	Quantity int
}
