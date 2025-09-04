package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// TelemetryClient handles communication with Telemetry Service
type TelemetryClient struct {
	baseURL    string
	httpClient *http.Client
}

// TelemetryReading represents a telemetry reading
type TelemetryReading struct {
	ID         int64     `json:"id"`
	DeviceID   string    `json:"device_id"`
	DeviceType string    `json:"device_type"`
	Location   string    `json:"location"`
	MetricName string    `json:"metric_name"`
	Value      float64   `json:"value"`
	Unit       string    `json:"unit"`
	Timestamp  time.Time `json:"timestamp"`
	RawData    *string   `json:"raw_data"`
	Quality    string    `json:"quality"`
	CreatedAt  time.Time `json:"created_at"`
}

// TelemetryReadingCreate represents a request to create a telemetry reading
type TelemetryReadingCreate struct {
	DeviceID   string    `json:"device_id"`
	DeviceType string    `json:"device_type"`
	Location   string    `json:"location"`
	MetricName string    `json:"metric_name"`
	Value      float64   `json:"value"`
	Unit       string    `json:"unit"`
	Timestamp  time.Time `json:"timestamp"`
	RawData    *string   `json:"raw_data,omitempty"`
	Quality    string    `json:"quality,omitempty"`
}

// TelemetryAlert represents an alert
type TelemetryAlert struct {
	ID               int64      `json:"id"`
	DeviceID         string     `json:"device_id"`
	DeviceType       string     `json:"device_type"`
	Location         string     `json:"location"`
	AlertType        string     `json:"alert_type"`
	Severity         string     `json:"severity"`
	CurrentValue     float64    `json:"current_value"`
	ThresholdValue   float64    `json:"threshold_value"`
	Message          string     `json:"message"`
	IsActive         string     `json:"is_active"`
	AcknowledgedAt   *time.Time `json:"acknowledged_at"`
	ResolvedAt       *time.Time `json:"resolved_at"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// TelemetryStats represents telemetry statistics
type TelemetryStats struct {
	TotalReadings        int64      `json:"total_readings"`
	TotalDevices         int64      `json:"total_devices"`
	ActiveAlerts         int64      `json:"active_alerts"`
	LastReadingTime      *time.Time `json:"last_reading_time"`
	ReadingsPerHour      float64    `json:"readings_per_hour"`
	AverageQualityScore  float64    `json:"average_quality_score"`
}

// LocationStats represents location-based statistics
type LocationStats struct {
	Location           string     `json:"location"`
	DeviceCount        int64      `json:"device_count"`
	TotalReadings      int64      `json:"total_readings"`
	LatestReadingTime  *time.Time `json:"latest_reading_time"`
	AvgTemperature     *float64   `json:"avg_temperature"`
	ActiveAlerts       int64      `json:"active_alerts"`
}

// NewTelemetryClient creates a new telemetry service client
func NewTelemetryClient() *TelemetryClient {
	baseURL := os.Getenv("TELEMETRY_SERVICE_URL")
	if baseURL == "" {
		baseURL = "http://telemetry-service:8083"
	}

	return &TelemetryClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// CreateReading creates a new telemetry reading
func (c *TelemetryClient) CreateReading(ctx context.Context, reading TelemetryReadingCreate) (*TelemetryReading, error) {
	jsonData, err := json.Marshal(reading)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/v1/telemetry/readings", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result TelemetryReading
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// GetReadings retrieves telemetry readings with optional filters
func (c *TelemetryClient) GetReadings(ctx context.Context, deviceID, location, metricName string, limit int) ([]TelemetryReading, error) {
	url := fmt.Sprintf("%s/api/v1/telemetry/readings", c.baseURL)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add query parameters
	q := req.URL.Query()
	if deviceID != "" {
		q.Add("device_id", deviceID)
	}
	if location != "" {
		q.Add("location", location)
	}
	if metricName != "" {
		q.Add("metric_name", metricName)
	}
	if limit > 0 {
		q.Add("limit", fmt.Sprintf("%d", limit))
	}
	req.URL.RawQuery = q.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var readings []TelemetryReading
	err = json.Unmarshal(body, &readings)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return readings, nil
}

// GetLatestReading retrieves the latest reading for a device and metric
func (c *TelemetryClient) GetLatestReading(ctx context.Context, deviceID, metricName string) (*TelemetryReading, error) {
	url := fmt.Sprintf("%s/api/v1/telemetry/readings/latest/%s/%s", c.baseURL, deviceID, metricName)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil // No readings found
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var reading TelemetryReading
	err = json.Unmarshal(body, &reading)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &reading, nil
}

// GetDeviceSummary retrieves summary data for a device
func (c *TelemetryClient) GetDeviceSummary(ctx context.Context, deviceID string, hours int) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/v1/telemetry/devices/%s/summary", c.baseURL, deviceID)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if hours > 0 {
		q := req.URL.Query()
		q.Add("hours", fmt.Sprintf("%d", hours))
		req.URL.RawQuery = q.Encode()
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var summary map[string]interface{}
	err = json.Unmarshal(body, &summary)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return summary, nil
}

// GetAlerts retrieves telemetry alerts
func (c *TelemetryClient) GetAlerts(ctx context.Context, activeOnly bool) ([]TelemetryAlert, error) {
	url := fmt.Sprintf("%s/api/v1/telemetry/alerts", c.baseURL)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if activeOnly {
		q := req.URL.Query()
		q.Add("active_only", "true")
		req.URL.RawQuery = q.Encode()
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var alerts []TelemetryAlert
	err = json.Unmarshal(body, &alerts)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return alerts, nil
}

// GetTelemetryStats retrieves overall telemetry statistics
func (c *TelemetryClient) GetTelemetryStats(ctx context.Context) (*TelemetryStats, error) {
	url := fmt.Sprintf("%s/api/v1/telemetry/stats", c.baseURL)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var stats TelemetryStats
	err = json.Unmarshal(body, &stats)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &stats, nil
}

// GetLocationStats retrieves location-based statistics
func (c *TelemetryClient) GetLocationStats(ctx context.Context) ([]LocationStats, error) {
	url := fmt.Sprintf("%s/api/v1/telemetry/stats/locations", c.baseURL)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var stats []LocationStats
	err = json.Unmarshal(body, &stats)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return stats, nil
}