package models

type Item struct {
	ID    uint `gorm:"primaryKey"`
	Name  string
	Price int
}
