package repositories

import (
	"avito-shop/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByUsername(username string) (*models.User, error)
	GetUserByID(id uint) (*models.User, error)
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
	GetUserBalance(userID uint) (float64, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) CreateUser(user *models.User) error {
	result := r.db.Create(user)
	return result.Error
}

func (r *userRepository) UpdateUser(user *models.User) error {
	result := r.db.Save(user)
	return result.Error
}

func (r *userRepository) GetUserBalance(userID uint) (float64, error) {
	var balance float64
	if err := r.db.Model(&models.User{}).Where("id = ?", userID).Pluck("balance", &balance).Error; err != nil {
		return 0, err
	}
	return balance, nil
}
