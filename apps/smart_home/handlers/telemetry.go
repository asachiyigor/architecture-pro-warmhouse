package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"smarthome/services"

	"github.com/gin-gonic/gin"
)

// TelemetryHandler handles telemetry-related requests via API Gateway
type TelemetryHandler struct {
	telemetryClient *services.TelemetryGatewayClient
}

// NewTelemetryHandler creates a new TelemetryHandler
func NewTelemetryHandler(telemetryClient *services.TelemetryGatewayClient) *TelemetryHandler {
	return &TelemetryHandler{
		telemetryClient: telemetryClient,
	}
}

// RegisterRoutes registers the telemetry routes
func (h *TelemetryHandler) RegisterRoutes(router *gin.RouterGroup) {
	telemetry := router.Group("/telemetry")
	{
		telemetry.GET("/readings", h.GetReadings)
		telemetry.POST("/readings", h.CreateReading)
		telemetry.GET("/readings/latest/:deviceId/:metricName", h.GetLatestReading)
		telemetry.GET("/stats", h.GetStats)
	}
}

// GetReadings handles GET /api/v1/telemetry/readings
func (h *TelemetryHandler) GetReadings(c *gin.Context) {
	ctx := context.Background()
	
	// Extract query parameters
	deviceID := c.Query("device_id")
	location := c.Query("location")
	metricName := c.Query("metric_name")
	limitStr := c.Query("limit")
	
	var limit int
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	readings, err := h.telemetryClient.GetReadings(ctx, deviceID, location, metricName, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, readings)
}

// CreateReading handles POST /api/v1/telemetry/readings
func (h *TelemetryHandler) CreateReading(c *gin.Context) {
	ctx := context.Background()
	
	var reading services.TelemetryReadingCreate
	if err := c.ShouldBindJSON(&reading); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set timestamp if not provided
	if reading.Timestamp.IsZero() {
		reading.Timestamp = time.Now()
	}

	result, err := h.telemetryClient.CreateReading(ctx, reading)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetLatestReading handles GET /api/v1/telemetry/readings/latest/:deviceId/:metricName
func (h *TelemetryHandler) GetLatestReading(c *gin.Context) {
	ctx := context.Background()
	
	deviceID := c.Param("deviceId")
	metricName := c.Param("metricName")
	
	if deviceID == "" || metricName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Device ID and metric name are required"})
		return
	}

	reading, err := h.telemetryClient.GetLatestReading(ctx, deviceID, metricName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	if reading == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No readings found"})
		return
	}

	c.JSON(http.StatusOK, reading)
}

// GetStats handles GET /api/v1/telemetry/stats
func (h *TelemetryHandler) GetStats(c *gin.Context) {
	ctx := context.Background()
	
	stats, err := h.telemetryClient.GetTelemetryStats(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}