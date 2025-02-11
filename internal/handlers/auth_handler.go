package handlers

import (
	"avito-shop/internal/middleware"
	"avito-shop/internal/models"
	"avito-shop/internal/repositories"
	"avito-shop/internal/services"
	"log"

	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userRepo *repositories.UserRepository
}

func NewAuthHandler(userRepo *repositories.UserRepository) *AuthHandler {
	return &AuthHandler{userRepo: userRepo}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Login request JSON parse error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := h.userRepo.GetUserByUsername(req.Username)
	if err != nil {
		log.Println("User not found:", req.Username, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Проверка пароля
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		log.Println("Password mismatch for user:", req.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Генерация токена (используем `Username`, а не `ID`)
	token, err := middleware.GenerateJWT(user.ID)
	if err != nil {
		log.Println("JWT generation failed:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	hashedPassword := services.HashPassword(req.Password)

	user := models.User{
		Username: req.Username,
		Password: hashedPassword,
		Balance:  1000,
	}

	if err := h.userRepo.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

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

	// Проверяем получателя
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

	// Обновляем балансы
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
