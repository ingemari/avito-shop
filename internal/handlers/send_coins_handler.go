package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *AuthHandler) Transaction(c *gin.Context) {
	// Получаем userID из JWT-токена (middleware сохраняет его в контекст)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req struct {
		To     string `json:"to" binding:"required"`
		Amount int    `json:"amount" binding:"gt=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Находим пользователя, который делает перевод (из токена, а не из тела запроса)
	fromUser, err := h.userRepo.GetUserByID(userID.(uint)) // userID из JWT
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sender not found"})
		return
	}

	toUser, err := h.userRepo.GetUserByUsername(req.To)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipient not found"})
		return
	}

	if fromUser.Username == toUser.Username {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot transfer to yourself"})
		return
	}

	if fromUser.Balance < req.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient funds"})
		return
	}

	fromUser.Balance -= req.Amount
	toUser.Balance += req.Amount

	// Сохраняем изменения
	if err := h.userRepo.UpdateUser(fromUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update sender balance"})
		return
	}

	if err := h.userRepo.UpdateUser(toUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update recipient balance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Transaction successful",
		"from":         fromUser.Username,
		"to":           toUser.Username,
		"amount":       req.Amount,
		"from_balance": fromUser.Balance,
		"to_balance":   toUser.Balance,
	})
}
