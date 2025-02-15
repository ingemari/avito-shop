package handlers

import (
	"avito-shop/internal/services"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	service services.TransactionService
}

func NewTransactionHandler(service services.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

func (h *TransactionHandler) SendCoins(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req struct {
		ToUser string `json:"toUser" binding:"required"`
		Amount int    `json:"amount" binding:"gt=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := h.service.TransferCoins(userID.(uint), req.ToUser, req.Amount)
	if err != nil {
		if errors.Is(err, errors.New("insufficient funds")) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient funds"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction successful"})
}
