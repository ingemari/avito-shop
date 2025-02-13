package handlers

import (
	"avito-shop/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PurchaseHandler struct {
	service *services.PurchaseService
}

func NewPurchaseHandler(service *services.PurchaseService) *PurchaseHandler {
	return &PurchaseHandler{service: service}
}

func (h *PurchaseHandler) Buy(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Получаем item из параметра пути
	item := c.Param("item")
	if item == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Item is required"})
		return
	}

	newBalance, err := h.service.PurchaseItem(userID.(uint), item)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Purchase successful",
		"new_balance": newBalance,
	})
}
