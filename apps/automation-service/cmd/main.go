// @title Automation Service API
// @version 1.0
// @description Микросервис для управления сценариями автоматизации умного дома
// @host localhost:8086
// @BasePath /api/v1
package main

import (
	"automation-service/internal/handlers"
	"automation-service/internal/repository"
	"automation-service/internal/service"
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/lib/pq"

	// Import generated docs
	_ "automation-service/docs"
)

func main() {
	// Get database URL from environment
	dbURL := getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5437/smarthome_automation?sslmode=disable")
	
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
	automationRepo := repository.NewAutomationRepository(db)
	automationService := service.NewAutomationService(automationRepo)
	automationHandler := handlers.NewAutomationHandler(automationService)

	// Initialize Gin router
	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "automation-service",
		})
	})

	// Swagger endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	api := router.Group("/api/v1")
	automationHandler.RegisterRoutes(api)

	// Start server
	port := getEnv("PORT", ":8086")
	log.Printf("Automation Service starting on %s", port)
	
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