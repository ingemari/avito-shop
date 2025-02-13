package main

import (
	"avito-shop/internal/handlers"
	"avito-shop/internal/middleware"
	"avito-shop/internal/repositories"
	"avito-shop/internal/services"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=postgres password=secret dbname=avito port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database")
	}

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)
	itemsRepo := repositories.NewItemsRepository(db)
	inventoryRepo := repositories.NewInventoryRepository(db)

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
