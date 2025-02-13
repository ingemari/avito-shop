// package handlers

// import (
// 	"avito-shop/internal/repositories"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// )

// type InventoryHandler struct {
// 	inventoryRepo *repositories.InventoryRepository
// }

// func NewInventoryHandler(inventoryRepo *repositories.InventoryRepository) *InventoryHandler {
// 	return &InventoryHandler{inventoryRepo: inventoryRepo}
// }

// func (h *InventoryHandler) Inventory(c *gin.Context) {
// 	// Получаем userID из JWT-токена (middleware передает его в `context`)
// 	userID, exists := c.Get("userID")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 		return
// 	}

// 	// Получаем инвентарь пользователя
// 	items, err := h.inventoryRepo.GetUserInventory(userID.(uint))
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get inventory"})
// 		return
// 	}

//		// Отправляем JSON-ответ
//		c.JSON(http.StatusOK, gin.H{"inventory": items})
//	}
package handlers

import (
	"avito-shop/internal/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
)

type InfoHandler struct {
	inventoryRepo   *repositories.InventoryRepository
	transactionRepo *repositories.TransactionRepository
	balanceService  *repositories.UserRepository
}

func NewInfoHandler(inventoryRepo *repositories.InventoryRepository, transactionRepo *repositories.TransactionRepository, balanceService *repositories.UserRepository) *InfoHandler {
	return &InfoHandler{
		inventoryRepo:   inventoryRepo,
		transactionRepo: transactionRepo,
		balanceService:  balanceService,
	}
}

func (h *InfoHandler) Info(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	balance, err := h.balanceService.GetUserBalance(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get balance"})
		return
	}

	items, err := h.inventoryRepo.GetUserInventory(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get inventory"})
		return
	}

	transactions, err := h.transactionRepo.GetUserTransactions(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get transactions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"balance":      balance,
		"inventory":    items,
		"transactions": transactions,
	})
}
