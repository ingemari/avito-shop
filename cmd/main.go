package main

import (
	"avito-shop/internal/db"
	"avito-shop/internal/handlers"
	"avito-shop/internal/middleware"
	"avito-shop/internal/repositories"
	"avito-shop/internal/services"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	db := db.NewDB()

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db.DB)
	transactionRepo := repositories.NewTransactionRepository(db.DB)
	itemsRepo := repositories.NewItemsRepository(db.DB)
	inventoryRepo := repositories.NewInventoryRepository(db.DB)

	// Initialize services
	authService := services.NewAuthService(userRepo)
	transactionService := services.NewTransactionService(userRepo, transactionRepo)
	purchaseService := services.NewPurchaseService(userRepo, itemsRepo, inventoryRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)
	purchaseHandler := handlers.NewPurchaseHandler(purchaseService)
	infoHandler := handlers.NewInfoHandler(inventoryRepo, transactionRepo, userRepo)

	r := gin.Default()

	// Public routes
	r.POST("/api/auth", authHandler.Login)

	// Protected routes with JWT middleware
	protected := r.Group("/api")
	protected.Use(middleware.JWTAuthMiddleware())
	protected.POST("/sendCoin", transactionHandler.SendCoins)
	protected.POST("/buy/:item", purchaseHandler.Buy)
	protected.GET("/info", infoHandler.Info)

	log.Println("Server is running on port 8080")
	r.Run(":8080")
}
