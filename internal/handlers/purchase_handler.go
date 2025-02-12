package handlers

import (
	"avito-shop/internal/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PurchaseHandler struct {
	userRepo      *repositories.UserRepository
	inventoryRepo *repositories.InventoryRepository
}

func NewPurchaseHandler(userRepo *repositories.UserRepository, inventoryRepo *repositories.InventoryRepository) *PurchaseHandler {
	return &PurchaseHandler{userRepo: userRepo, inventoryRepo: inventoryRepo}
}

var Items = map[string]int{
	"t-shirt":    80,
	"cup":        20,
	"book":       50,
	"pen":        10,
	"powerbank":  200,
	"hoody":      300,
	"umbrella":   200,
	"socks":      10,
	"wallet":     50,
	"pink-hoody": 500,
}

func (h *PurchaseHandler) Purchase(c *gin.Context) {
	// Получаем userID из JWT-токена (middleware сохраняет его в контекст)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req struct {
		Item     string `json:"item" binding:"required"`
		Quantity int    `json:"quantity" binding:"gt=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Проверяем, есть ли товар в магазине
	price, itemExists := Items[req.Item]
	if !itemExists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Item not found"})
		return
	}

	// Получаем данные пользователя
	user, err := h.userRepo.GetUserByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Проверяем, хватает ли средств
	totalCost := price * req.Quantity
	if user.Balance < totalCost {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient funds"})
		return
	}

	// Списываем деньги
	user.Balance -= totalCost

	// Обновляем баланс пользователя
	if err := h.userRepo.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update balance"})
		return
	}

	// Добавляем товар в инвентарь
	if err := h.inventoryRepo.AddItem(user.ID, req.Item, req.Quantity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add item to inventory"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Purchase successful", "new_balance": user.Balance})
}
