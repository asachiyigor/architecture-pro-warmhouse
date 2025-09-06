package gateway

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Config holds the configuration for the API Gateway
type Config struct {
	SmartHomeURL           string
	DeviceServiceURL       string
	TelemetryServiceURL    string
	UserServiceURL         string
	BillingServiceURL      string
	AutomationServiceURL   string
	NotificationServiceURL string
}

// Gateway represents the API Gateway
type Gateway struct {
	config Config
	client *http.Client
}

// New creates a new Gateway instance
func New(config Config) *Gateway {
	return &Gateway{
		config: config,
		client: &http.Client{},
	}
}

// RegisterRoutes registers all gateway routes
func (gw *Gateway) RegisterRoutes(router *gin.Engine) {
	// Gateway info endpoints (for Swagger documentation)
	router.GET("/gateway/status", gw.getGatewayStatus)
	router.GET("/gateway/services", gw.getServices)

	api := router.Group("/api/v1")

	// Phase 1: Everything goes to Smart Home App (current monolith)
	// This provides a transparent proxy while we migrate services

	// Smart Home App routes (sensors and legacy endpoints)
	api.Any("/sensors/*path", gw.proxyToSmartHome)

	// Migrated microservice routes (Phase 2)
	api.Any("/telemetry/*path", gw.proxyToTelemetryService) // ✅ MIGRATED to Telemetry Service

	// Future microservice routes (still proxy to Smart Home App)
	api.Any("/devices/*path", gw.proxyToSmartHome)                 // Will migrate to Device Service
	api.Any("/users/*path", gw.proxyToSmartHome)                   // Will migrate to User Service
	api.Any("/billing/*path", gw.proxyToSmartHome)                 // Will migrate to Billing Service
	api.Any("/automation/*path", gw.proxyToSmartHome)              // Will migrate to Automation Service
	api.Any("/notifications/*path", gw.proxyToNotificationService) // Direct to Notification Service

	// Catch-all route for any other paths
	router.NoRoute(gw.proxyToSmartHome)
}

// proxyToSmartHome proxies requests to the Smart Home App
func (gw *Gateway) proxyToSmartHome(c *gin.Context) {
	gw.proxyRequest(c, gw.config.SmartHomeURL)
}

// proxyToDeviceService proxies requests to the Device Service
func (gw *Gateway) proxyToDeviceService(c *gin.Context) {
	gw.proxyRequest(c, gw.config.DeviceServiceURL)
}

// proxyToTelemetryService proxies requests to the Telemetry Service
func (gw *Gateway) proxyToTelemetryService(c *gin.Context) {
	gw.proxyRequest(c, gw.config.TelemetryServiceURL)
}

// proxyToUserService proxies requests to the User Service
func (gw *Gateway) proxyToUserService(c *gin.Context) {
	gw.proxyRequest(c, gw.config.UserServiceURL)
}

// proxyToBillingService proxies requests to the Billing Service
func (gw *Gateway) proxyToBillingService(c *gin.Context) {
	gw.proxyRequest(c, gw.config.BillingServiceURL)
}

// proxyToAutomationService proxies requests to the Automation Service
func (gw *Gateway) proxyToAutomationService(c *gin.Context) {
	gw.proxyRequest(c, gw.config.AutomationServiceURL)
}

// proxyToNotificationService proxies requests to the Notification Service
func (gw *Gateway) proxyToNotificationService(c *gin.Context) {
	gw.proxyRequest(c, gw.config.NotificationServiceURL)
}

// proxyRequest performs the actual request proxying
func (gw *Gateway) proxyRequest(c *gin.Context, targetURL string) {
	// Parse target URL
	target, err := url.Parse(targetURL)
	if err != nil {
		log.Printf("Error parsing target URL %s: %v", targetURL, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid target URL"})
		return
	}

	// Build the target request URL
	targetReqURL := target.ResolveReference(&url.URL{
		Path:     c.Request.URL.Path,
		RawQuery: c.Request.URL.RawQuery,
	})

	// Read request body
	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, err = io.ReadAll(c.Request.Body)
		if err != nil {
			log.Printf("Error reading request body: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading request body"})
			return
		}
	}

	// Create new request
	proxyReq, err := http.NewRequest(c.Request.Method, targetReqURL.String(), bytes.NewReader(bodyBytes))
	if err != nil {
		log.Printf("Error creating proxy request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating proxy request"})
		return
	}

	// Copy headers (excluding hop-by-hop headers)
	for key, values := range c.Request.Header {
		if !isHopByHopHeader(key) {
			for _, value := range values {
				proxyReq.Header.Add(key, value)
			}
		}
	}

	// Set X-Forwarded headers
	proxyReq.Header.Set("X-Forwarded-For", c.ClientIP())
	proxyReq.Header.Set("X-Forwarded-Host", c.Request.Host)
	proxyReq.Header.Set("X-Forwarded-Proto", getScheme(c.Request))

	// Log the proxy request
	log.Printf("Proxying %s %s to %s", c.Request.Method, c.Request.URL.Path, targetReqURL.String())

	// Make the request
	resp, err := gw.client.Do(proxyReq)
	if err != nil {
		log.Printf("Error making proxy request to %s: %v", targetReqURL.String(), err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "Service unavailable"})
		return
	}
	defer resp.Body.Close()

	// Copy response headers
	for key, values := range resp.Header {
		if !isHopByHopHeader(key) {
			for _, value := range values {
				c.Header(key, value)
			}
		}
	}

	// Copy response body
	c.Status(resp.StatusCode)
	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		log.Printf("Error copying response body: %v", err)
	}
}

// isHopByHopHeader checks if a header is hop-by-hop
func isHopByHopHeader(header string) bool {
	hopByHopHeaders := []string{
		"Connection",
		"Keep-Alive",
		"Proxy-Authenticate",
		"Proxy-Authorization",
		"Te",
		"Trailers",
		"Transfer-Encoding",
		"Upgrade",
	}

	header = strings.ToLower(header)
	for _, hopHeader := range hopByHopHeaders {
		if strings.ToLower(hopHeader) == header {
			return true
		}
	}
	return false
}

// getScheme returns the request scheme
func getScheme(req *http.Request) string {
	if req.TLS != nil {
		return "https"
	}
	return "http"
}

// GatewayStatus represents the gateway status response
type GatewayStatus struct {
	Status    string            `json:"status" example:"running"`
	Version   string            `json:"version" example:"1.0.0"`
	Services  map[string]string `json:"services"`
	Timestamp int64             `json:"timestamp" example:"1693123456"`
}

// ServiceInfo represents service information
type ServiceInfo struct {
	Name        string `json:"name" example:"user-service"`
	URL         string `json:"url" example:"http://user-service:8084"`
	Description string `json:"description" example:"Управление пользователями и домохозяйствами"`
}

// getGatewayStatus returns the current status of the API Gateway
// @Summary Get API Gateway Status
// @Description Получить статус и информацию о API Gateway
// @Tags Gateway
// @Produce json
// @Success 200 {object} GatewayStatus
// @Router /gateway/status [get]
func (gw *Gateway) getGatewayStatus(c *gin.Context) {
	status := GatewayStatus{
		Status:  "running",
		Version: "1.0.0",
		Services: map[string]string{
			"smart-home":   gw.config.SmartHomeURL,
			"device":       gw.config.DeviceServiceURL,
			"telemetry":    gw.config.TelemetryServiceURL,
			"user":         gw.config.UserServiceURL,
			"billing":      gw.config.BillingServiceURL,
			"automation":   gw.config.AutomationServiceURL,
			"notification": gw.config.NotificationServiceURL,
		},
		Timestamp: time.Now().Unix(),
	}
	c.JSON(http.StatusOK, status)
}

// getServices returns information about available services
// @Summary Get Available Services
// @Description Получить список доступных микросервисов
// @Tags Gateway
// @Produce json
// @Success 200 {array} ServiceInfo
// @Router /gateway/services [get]
func (gw *Gateway) getServices(c *gin.Context) {
	services := []ServiceInfo{
		{
			Name:        "smart-home-app",
			URL:         gw.config.SmartHomeURL,
			Description: "Монолитное приложение Smart Home (legacy)",
		},
		{
			Name:        "device-service",
			URL:         gw.config.DeviceServiceURL,
			Description: "Микросервис управления устройствами",
		},
		{
			Name:        "telemetry-service",
			URL:         gw.config.TelemetryServiceURL,
			Description: "Микросервис сбора телеметрии",
		},
		{
			Name:        "user-service",
			URL:         gw.config.UserServiceURL,
			Description: "Микросервис управления пользователями",
		},
		{
			Name:        "billing-service",
			URL:         gw.config.BillingServiceURL,
			Description: "Микросервис управления платежами",
		},
		{
			Name:        "automation-service",
			URL:         gw.config.AutomationServiceURL,
			Description: "Микросервис автоматизации",
		},
		{
			Name:        "notification-service",
			URL:         gw.config.NotificationServiceURL,
			Description: "Микросервис уведомлений",
		},
	}
	c.JSON(http.StatusOK, services)
}
