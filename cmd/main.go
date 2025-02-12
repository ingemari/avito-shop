package main

import (
	"avito-shop/internal/handlers"
	"avito-shop/internal/middleware"
	"avito-shop/internal/repositories"
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

	userRepo := repositories.NewUserRepository(db)
	authHandler := handlers.NewAuthHandler(userRepo)
	inventoryRepo := repositories.NewInventoryRepository(db)
	purchaseHandler := handlers.NewPurchaseHandler(userRepo, inventoryRepo)
	inventoryHandler := handlers.NewInventoryHandler(inventoryRepo)

	r := gin.Default()

	// Без аутентификации
	r.POST("/api/auth", authHandler.Login)

	// Все защищенные маршруты
	protected := r.Group("/api")
	protected.Use(middleware.JWTAuthMiddleware()) // JWT-мидлвар
	protected.GET("/inventory", inventoryHandler.Inventory)
	protected.POST("/transaction", authHandler.Transaction)
	protected.POST("/purchase", purchaseHandler.Purchase)

	log.Println("Server is running on port 8080")
	r.Run(":8080")
}
