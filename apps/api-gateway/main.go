// @title Smart Home API Gateway
// @version 1.0
// @description АПИ шлюз для микросервисной архитектуры Smart Home Pro
// @host localhost:8000
// @BasePath /
package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"api-gateway/internal/gateway"
	"api-gateway/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// Import generated docs
	_ "api-gateway/docs"
)

func main() {
	// Initialize Gin router
	router := gin.Default()

	// CORS configuration
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	router.Use(cors.New(config))

	// Add logging middleware
	router.Use(middleware.Logger())

	// Swagger endpoint - aggregated documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Add health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"service": "api-gateway",
			"timestamp": time.Now().Unix(),
		})
	})

	// Initialize gateway with service URLs
	gw := gateway.New(gateway.Config{
		SmartHomeURL:           getEnv("SMART_HOME_URL", "http://app:8080"),
		DeviceServiceURL:       getEnv("DEVICE_SERVICE_URL", "http://device-service:8082"),
		TelemetryServiceURL:    getEnv("TELEMETRY_SERVICE_URL", "http://telemetry-service:8083"),
		UserServiceURL:         getEnv("USER_SERVICE_URL", "http://user-service:8084"),
		BillingServiceURL:      getEnv("BILLING_SERVICE_URL", "http://billing-service:8085"),
		AutomationServiceURL:   getEnv("AUTOMATION_SERVICE_URL", "http://automation-service:8087"),
		NotificationServiceURL: getEnv("NOTIFICATION_SERVICE_URL", "http://notification-service:8086"),
	})

	// Register routes
	gw.RegisterRoutes(router)

	// Start server
	port := getEnv("PORT", ":8000")
	log.Printf("API Gateway starting on port %s", port)
	
	srv := &http.Server{
		Addr:    port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down API Gateway...")

	log.Println("API Gateway stopped")
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}