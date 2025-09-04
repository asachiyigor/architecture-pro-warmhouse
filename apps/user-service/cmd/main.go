// @title User Service API
// @version 1.0
// @description Микросервис для управления пользователями и домохозяйствами
// @host localhost:8084
// @BasePath /api/v1
package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"user-service/internal/handlers"
	"user-service/internal/repository"
	"user-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/lib/pq"

	// Import generated docs
	_ "user-service/docs"
)

func main() {
	// Get database URL from environment
	dbURL := getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5435/smarthome_users?sslmode=disable")
	
	// Connect to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Connected to database successfully")

	// Initialize layers
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	// Initialize Gin router
	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "user-service",
		})
	})

	// Swagger endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	api := router.Group("/api/v1")
	userHandler.RegisterRoutes(api)

	// Start server
	port := getEnv("PORT", ":8084")
	log.Printf("User Service starting on %s", port)
	
	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}