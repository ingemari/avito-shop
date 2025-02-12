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

	// Проверяем, существует ли пользователь
	user, err := h.userRepo.GetUserByUsername(req.Username)
	if err != nil {
		log.Println("User not found, creating new:", req.Username)

		// Создаем нового пользователя
		hashedPassword := services.HashPassword(req.Password)
		newUser := models.User{
			Username: req.Username,
			Password: hashedPassword,
			Balance:  1000, // Стартовый баланс
		}

		if err := h.userRepo.CreateUser(&newUser); err != nil {
			log.Println("Failed to create user:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		user = &newUser
	} else {
		// Если пользователь найден, проверяем пароль
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			log.Println("Password mismatch for user:", req.Username)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
	}

	// Генерация JWT-токена
	token, err := middleware.GenerateJWT(user.ID)
	if err != nil {
		log.Println("JWT generation failed:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
