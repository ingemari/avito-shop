package models

type Transaction struct {
	ID       uint `gorm:"primaryKey"`
	FromUser uint `gorm:"index"`
	ToUser   uint `gorm:"index"`
	Amount   int
}
