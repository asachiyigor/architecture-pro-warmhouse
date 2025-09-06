package service

import (
	"billing-service/internal/models"
	"billing-service/internal/repository"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type BillingService struct {
	billingRepo *repository.BillingRepository
}

func NewBillingService(billingRepo *repository.BillingRepository) *BillingService {
	return &BillingService{billingRepo: billingRepo}
}

// GetPlans retrieves all available pricing plans
func (s *BillingService) GetPlans() ([]*models.PricingPlan, error) {
	return s.billingRepo.GetAllPlans()
}

// GetPlan retrieves a specific pricing plan
func (s *BillingService) GetPlan(id string) (*models.PricingPlan, error) {
	return s.billingRepo.GetPlanByID(id)
}

// CreateSubscription creates a new subscription for a user
func (s *BillingService) CreateSubscription(req models.CreateSubscriptionRequest) (*models.Subscription, error) {
	// Validate that the plan exists
	plan, err := s.billingRepo.GetPlanByID(req.PlanID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("pricing plan not found")
		}
		return nil, err
	}
	
	if !plan.IsActive {
		return nil, errors.New("pricing plan is not active")
	}
	
	// Check if user already has an active subscription
	userSubs, err := s.billingRepo.GetUserSubscriptions(req.UserID)
	if err != nil {
		return nil, err
	}
	
	for _, sub := range userSubs {
		if sub.IsActive() {
			return nil, errors.New("user already has an active subscription")
		}
	}
	
	// Create new subscription
	subscription := models.NewSubscription(req.UserID, req.PlanID)
	
	// Set end date based on billing interval
	if plan.BillingInterval == "monthly" {
		endDate := subscription.StartedAt.AddDate(0, 1, 0)
		subscription.EndsAt = &endDate
	} else if plan.BillingInterval == "yearly" {
		endDate := subscription.StartedAt.AddDate(1, 0, 0)
		subscription.EndsAt = &endDate
	}
	
	err = s.billingRepo.CreateSubscription(subscription)
	if err != nil {
		return nil, err
	}
	
	// Attach plan details
	subscription.Plan = plan
	
	return subscription, nil
}

// GetSubscription retrieves a subscription by ID
func (s *BillingService) GetSubscription(id string) (*models.Subscription, error) {
	return s.billingRepo.GetSubscriptionByID(id)
}

// GetUserSubscriptions retrieves all subscriptions for a user
func (s *BillingService) GetUserSubscriptions(userID string) ([]*models.Subscription, error) {
	return s.billingRepo.GetUserSubscriptions(userID)
}

// CancelSubscription cancels a user's subscription
func (s *BillingService) CancelSubscription(subscriptionID string) error {
	subscription, err := s.billingRepo.GetSubscriptionByID(subscriptionID)
	if err != nil {
		return err
	}
	
	if subscription.Status != "active" {
		return errors.New("subscription is not active")
	}
	
	return s.billingRepo.UpdateSubscriptionStatus(subscriptionID, "cancelled")
}

// ProcessPayment processes a payment for a subscription
func (s *BillingService) ProcessPayment(req models.PaymentRequest) (*models.Payment, error) {
	// Validate subscription exists and is active
	subscription, err := s.billingRepo.GetSubscriptionByID(req.SubscriptionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("subscription not found")
		}
		return nil, err
	}
	
	if !subscription.IsActive() {
		return nil, errors.New("subscription is not active")
	}
	
	// Create payment record
	payment := models.NewPayment(req.SubscriptionID, req.Amount, req.PaymentMethod)
	
	err = s.billingRepo.CreatePayment(payment)
	if err != nil {
		return nil, err
	}
	
	// Simulate payment processing
	// In real implementation, this would integrate with payment gateway
	err = s.processPaymentWithGateway(payment)
	if err != nil {
		// Update payment status to failed
		s.billingRepo.UpdatePaymentStatus(payment.ID, "failed", "")
		return nil, fmt.Errorf("payment processing failed: %v", err)
	}
	
	// Update payment status to completed
	externalID := fmt.Sprintf("ext_%s", payment.ID[:8])
	err = s.billingRepo.UpdatePaymentStatus(payment.ID, "completed", externalID)
	if err != nil {
		return nil, err
	}
	
	payment.Status = "completed"
	payment.ExternalID = externalID
	processedAt := time.Now()
	payment.ProcessedAt = &processedAt
	
	return payment, nil
}

// GetPaymentHistory retrieves payment history for a subscription
func (s *BillingService) GetPaymentHistory(subscriptionID string) ([]*models.Payment, error) {
	return s.billingRepo.GetPaymentsBySubscription(subscriptionID)
}

// GetBillingStats retrieves billing statistics
func (s *BillingService) GetBillingStats() (*models.BillingStats, error) {
	return s.billingRepo.GetBillingStats()
}

// processPaymentWithGateway simulates payment gateway integration
func (s *BillingService) processPaymentWithGateway(payment *models.Payment) error {
	// Simulate payment processing delay
	time.Sleep(100 * time.Millisecond)
	
	// Simulate success (90% success rate)
	// In real implementation, this would call actual payment gateway
	if payment.Amount <= 0 {
		return errors.New("invalid payment amount")
	}
	
	// Simulate some failures for testing
	if payment.PaymentMethod == "invalid_card" {
		return errors.New("invalid payment method")
	}
	
	return nil
}