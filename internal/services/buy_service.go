package services

import (
	"avito-shop/internal/repositories"
	"errors"
)

type PurchaseService struct {
	userRepo  *repositories.UserRepository
	itemsRepo *repositories.ItemsRepository
	invRepo   *repositories.InventoryRepository
}

func NewPurchaseService(userRepo *repositories.UserRepository, itemsRepo *repositories.ItemsRepository, invRepo *repositories.InventoryRepository) *PurchaseService {
	return &PurchaseService{userRepo: userRepo, itemsRepo: itemsRepo, invRepo: invRepo}
}

func (s *PurchaseService) PurchaseItem(userID uint, itemName string) (int, error) {
	item, err := s.itemsRepo.GetItemByName(itemName)
	if err != nil {
		return 0, errors.New("item not found")
	}

	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return 0, errors.New("user not found")
	}

	totalCost := item.Price
	if user.Balance < totalCost {
		return user.Balance, errors.New("insufficient funds")
	}

	user.Balance -= totalCost
	if err := s.userRepo.UpdateUser(user); err != nil {
		return user.Balance, errors.New("failed to update balance")
	}

	if err := s.invRepo.AddItem(user.ID, item.Name, 1); err != nil {
		return user.Balance, errors.New("failed to add item to inventory")
	}

	return user.Balance, nil
}
