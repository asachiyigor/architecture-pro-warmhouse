package models

import (
	"time"
	"github.com/google/uuid"
)

// PricingPlan represents a subscription plan
type PricingPlan struct {
	ID              string  `json:"id" db:"id"`
	Name            string  `json:"name" db:"name"`
	Description     string  `json:"description" db:"description"`
	Price           float64 `json:"price" db:"price"`
	BillingInterval string  `json:"billing_interval" db:"billing_interval"` // monthly, yearly
	DeviceLimit     int     `json:"device_limit" db:"device_limit"`
	IsActive        bool    `json:"is_active" db:"is_active"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// Subscription represents a user's subscription
type Subscription struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	PlanID    string    `json:"plan_id" db:"plan_id"`
	Status    string    `json:"status" db:"status"` // active, cancelled, expired, suspended
	StartedAt time.Time `json:"started_at" db:"started_at"`
	EndsAt    *time.Time `json:"ends_at" db:"ends_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	
	// Embedded plan details for convenience
	Plan *PricingPlan `json:"plan,omitempty"`
}

// Payment represents a payment record
type Payment struct {
	ID             string    `json:"id" db:"id"`
	SubscriptionID string    `json:"subscription_id" db:"subscription_id"`
	Amount         float64   `json:"amount" db:"amount"`
	Currency       string    `json:"currency" db:"currency"`
	Status         string    `json:"status" db:"status"` // pending, completed, failed, refunded
	PaymentMethod  string    `json:"payment_method" db:"payment_method"`
	ExternalID     string    `json:"external_id" db:"external_id"` // Payment gateway transaction ID
	ProcessedAt    *time.Time `json:"processed_at" db:"processed_at"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

// CreateSubscriptionRequest represents request to create subscription
type CreateSubscriptionRequest struct {
	UserID string `json:"user_id" binding:"required"`
	PlanID string `json:"plan_id" binding:"required"`
}

// PaymentRequest represents request to process payment
type PaymentRequest struct {
	SubscriptionID string  `json:"subscription_id" binding:"required"`
	Amount         float64 `json:"amount" binding:"required"`
	PaymentMethod  string  `json:"payment_method" binding:"required"`
}

// BillingStats represents billing statistics
type BillingStats struct {
	TotalActiveSubscriptions int     `json:"total_active_subscriptions"`
	TotalRevenue            float64 `json:"total_revenue"`
	MonthlyRevenue          float64 `json:"monthly_revenue"`
	TotalPayments           int     `json:"total_payments"`
}

// NewSubscription creates a new subscription
func NewSubscription(userID, planID string) *Subscription {
	now := time.Now()
	return &Subscription{
		ID:        uuid.New().String(),
		UserID:    userID,
		PlanID:    planID,
		Status:    "active",
		StartedAt: now,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// NewPayment creates a new payment record
func NewPayment(subscriptionID string, amount float64, paymentMethod string) *Payment {
	return &Payment{
		ID:             uuid.New().String(),
		SubscriptionID: subscriptionID,
		Amount:         amount,
		Currency:       "USD",
		Status:         "pending",
		PaymentMethod:  paymentMethod,
		CreatedAt:      time.Now(),
	}
}

// IsActive checks if subscription is currently active
func (s *Subscription) IsActive() bool {
	return s.Status == "active" && (s.EndsAt == nil || s.EndsAt.After(time.Now()))
}