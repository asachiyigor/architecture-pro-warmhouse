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

// DeviceClient handles communication with Device Management Service
type DeviceClient struct {
	baseURL    string
	httpClient *http.Client
}

// Device represents a device from the Device Management Service
type Device struct {
	ID               int64     `json:"id"`
	Name             string    `json:"name"`
	Type             string    `json:"type"`
	Location         string    `json:"location"`
	Status           string    `json:"status"`
	MacAddress       *string   `json:"macAddress"`
	IPAddress        *string   `json:"ipAddress"`
	FirmwareVersion  *string   `json:"firmwareVersion"`
	UserID           *int64    `json:"userId"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        *time.Time `json:"updatedAt"`
	LastSeen         *time.Time `json:"lastSeen"`
}

// DeviceCreateRequest represents a request to create a device
type DeviceCreateRequest struct {
	Name     string  `json:"name"`
	Type     string  `json:"type"`
	Location string  `json:"location"`
	UserID   *int64  `json:"userId,omitempty"`
}

// DeviceStats represents device statistics
type DeviceStats struct {
	TotalDevices   int64 `json:"totalDevices"`
	OnlineDevices  int64 `json:"onlineDevices"`
	OfflineDevices int64 `json:"offlineDevices"`
	ErrorDevices   int64 `json:"errorDevices"`
}

// NewDeviceClient creates a new device service client
func NewDeviceClient() *DeviceClient {
	baseURL := os.Getenv("DEVICE_SERVICE_URL")
	if baseURL == "" {
		baseURL = "http://device-service:8082"
	}

	return &DeviceClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetDevices retrieves all devices
func (c *DeviceClient) GetDevices(ctx context.Context) ([]Device, error) {
	url := fmt.Sprintf("%s/api/v1/devices", c.baseURL)
	
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

	var devices []Device
	err = json.Unmarshal(body, &devices)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return devices, nil
}

// GetDeviceByID retrieves a specific device by ID
func (c *DeviceClient) GetDeviceByID(ctx context.Context, deviceID int64) (*Device, error) {
	url := fmt.Sprintf("%s/api/v1/devices/%d", c.baseURL, deviceID)
	
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
		return nil, nil // Device not found
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var device Device
	err = json.Unmarshal(body, &device)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &device, nil
}

// CreateDevice creates a new device
func (c *DeviceClient) CreateDevice(ctx context.Context, deviceReq DeviceCreateRequest) (*Device, error) {
	jsonData, err := json.Marshal(deviceReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/v1/devices", c.baseURL)
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

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var device Device
	err = json.Unmarshal(body, &device)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &device, nil
}

// UpdateDeviceStatus updates the status of a device
func (c *DeviceClient) UpdateDeviceStatus(ctx context.Context, deviceID int64, status string) error {
	statusUpdate := map[string]string{"status": status}
	jsonData, err := json.Marshal(statusUpdate)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/v1/devices/%d/status", c.baseURL, deviceID)
	req, err := http.NewRequestWithContext(ctx, "PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

// GetDeviceStats retrieves device statistics
func (c *DeviceClient) GetDeviceStats(ctx context.Context) (*DeviceStats, error) {
	url := fmt.Sprintf("%s/api/v1/devices/stats", c.baseURL)
	
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

	var stats DeviceStats
	err = json.Unmarshal(body, &stats)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &stats, nil
}

// SendHeartbeat sends a heartbeat for a device
func (c *DeviceClient) SendHeartbeat(ctx context.Context, deviceID int64, ipAddress string) error {
	heartbeatData := map[string]string{"ipAddress": ipAddress}
	jsonData, err := json.Marshal(heartbeatData)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/v1/devices/%d/heartbeat", c.baseURL, deviceID)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}