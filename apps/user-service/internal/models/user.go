package models

import (
	"time"
	"github.com/google/uuid"
)

// User represents a system user
type User struct {
	ID        string    `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	FirstName string    `json:"first_name" db:"first_name"`
	LastName  string    `json:"last_name" db:"last_name"`
	Phone     string    `json:"phone" db:"phone"`
	Timezone  string    `json:"timezone" db:"timezone"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Household represents a user's household
type Household struct {
	ID          string     `json:"id" db:"id"`
	OwnerID     string     `json:"owner_id" db:"owner_id"`
	Name        string     `json:"name" db:"name"`
	Address     string     `json:"address" db:"address"`
	Timezone    string     `json:"timezone" db:"timezone"`
	Coordinates *string    `json:"coordinates" db:"coordinates"`
	WiFiSSID    string     `json:"wifi_ssid" db:"wifi_ssid"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// HouseholdMember represents membership in a household
type HouseholdMember struct {
	ID          string    `json:"id" db:"id"`
	HouseholdID string    `json:"household_id" db:"household_id"`
	UserID      string    `json:"user_id" db:"user_id"`
	Role        string    `json:"role" db:"role"`
	Permissions string    `json:"permissions" db:"permissions"` // JSON string
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// CreateUserRequest represents request to create a user
type CreateUserRequest struct {
	Email     string `json:"email" binding:"required,email"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Phone     string `json:"phone"`
	Timezone  string `json:"timezone"`
}

// CreateHouseholdRequest represents request to create a household
type CreateHouseholdRequest struct {
	Name        string `json:"name" binding:"required"`
	Address     string `json:"address" binding:"required"`
	Timezone    string `json:"timezone"`
	WiFiSSID    string `json:"wifi_ssid"`
}

// NewUser creates a new user with generated ID
func NewUser(req CreateUserRequest) *User {
	return &User{
		ID:        uuid.New().String(),
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Timezone:  req.Timezone,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// NewHousehold creates a new household with generated ID
func NewHousehold(ownerID string, req CreateHouseholdRequest) *Household {
	return &Household{
		ID:        uuid.New().String(),
		OwnerID:   ownerID,
		Name:      req.Name,
		Address:   req.Address,
		Timezone:  req.Timezone,
		WiFiSSID:  req.WiFiSSID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}