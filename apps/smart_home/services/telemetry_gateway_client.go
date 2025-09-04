package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// TelemetryGatewayClient handles communication with Telemetry Service via API Gateway
type TelemetryGatewayClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewTelemetryClientViaGateway creates a new telemetry client that uses API Gateway
func NewTelemetryClientViaGateway(apiGatewayURL string) *TelemetryGatewayClient {
	return &TelemetryGatewayClient{
		baseURL: apiGatewayURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// CreateReading creates a new telemetry reading via API Gateway
func (c *TelemetryGatewayClient) CreateReading(ctx context.Context, reading TelemetryReadingCreate) (*TelemetryReading, error) {
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
		return nil, fmt.Errorf("failed to make request via API Gateway: %w", err)
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

// GetReadings retrieves telemetry readings via API Gateway
func (c *TelemetryGatewayClient) GetReadings(ctx context.Context, deviceID, location, metricName string, limit int) ([]TelemetryReading, error) {
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
		return nil, fmt.Errorf("failed to make request via API Gateway: %w", err)
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

// GetLatestReading retrieves the latest reading for a device and metric via API Gateway
func (c *TelemetryGatewayClient) GetLatestReading(ctx context.Context, deviceID, metricName string) (*TelemetryReading, error) {
	url := fmt.Sprintf("%s/api/v1/telemetry/readings/latest/%s/%s", c.baseURL, deviceID, metricName)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request via API Gateway: %w", err)
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

// GetTelemetryStats retrieves overall telemetry statistics via API Gateway
func (c *TelemetryGatewayClient) GetTelemetryStats(ctx context.Context) (*TelemetryStats, error) {
	url := fmt.Sprintf("%s/api/v1/telemetry/stats", c.baseURL)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request via API Gateway: %w", err)
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