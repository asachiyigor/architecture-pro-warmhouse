package repository

import (
	"database/sql"
	"billing-service/internal/models"
	"time"
)

type BillingRepository struct {
	db *sql.DB
}

func NewBillingRepository(db *sql.DB) *BillingRepository {
	return &BillingRepository{db: db}
}

// GetAllPlans retrieves all active pricing plans
func (r *BillingRepository) GetAllPlans() ([]*models.PricingPlan, error) {
	query := `
		SELECT id, name, description, price, billing_interval, device_limit, is_active, created_at, updated_at
		FROM pricing_plans WHERE is_active = true ORDER BY price ASC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []*models.PricingPlan
	for rows.Next() {
		plan := &models.PricingPlan{}
		err := rows.Scan(&plan.ID, &plan.Name, &plan.Description, &plan.Price,
			&plan.BillingInterval, &plan.DeviceLimit, &plan.IsActive,
			&plan.CreatedAt, &plan.UpdatedAt)
		if err != nil {
			return nil, err
		}
		plans = append(plans, plan)
	}
	return plans, nil
}

// GetPlanByID retrieves a pricing plan by ID
func (r *BillingRepository) GetPlanByID(id string) (*models.PricingPlan, error) {
	query := `
		SELECT id, name, description, price, billing_interval, device_limit, is_active, created_at, updated_at
		FROM pricing_plans WHERE id = $1
	`
	plan := &models.PricingPlan{}
	err := r.db.QueryRow(query, id).Scan(&plan.ID, &plan.Name, &plan.Description, &plan.Price,
		&plan.BillingInterval, &plan.DeviceLimit, &plan.IsActive,
		&plan.CreatedAt, &plan.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return plan, nil
}

// CreateSubscription creates a new subscription
func (r *BillingRepository) CreateSubscription(subscription *models.Subscription) error {
	query := `
		INSERT INTO subscriptions (id, user_id, plan_id, status, started_at, ends_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.Exec(query, subscription.ID, subscription.UserID, subscription.PlanID,
		subscription.Status, subscription.StartedAt, subscription.EndsAt,
		subscription.CreatedAt, subscription.UpdatedAt)
	return err
}

// GetSubscriptionByID retrieves a subscription by ID
func (r *BillingRepository) GetSubscriptionByID(id string) (*models.Subscription, error) {
	query := `
		SELECT s.id, s.user_id, s.plan_id, s.status, s.started_at, s.ends_at, s.created_at, s.updated_at,
		       p.name, p.description, p.price, p.billing_interval, p.device_limit, p.is_active, p.created_at, p.updated_at
		FROM subscriptions s
		JOIN pricing_plans p ON s.plan_id = p.id
		WHERE s.id = $1
	`
	subscription := &models.Subscription{Plan: &models.PricingPlan{}}
	err := r.db.QueryRow(query, id).Scan(
		&subscription.ID, &subscription.UserID, &subscription.PlanID, &subscription.Status,
		&subscription.StartedAt, &subscription.EndsAt, &subscription.CreatedAt, &subscription.UpdatedAt,
		&subscription.Plan.Name, &subscription.Plan.Description, &subscription.Plan.Price,
		&subscription.Plan.BillingInterval, &subscription.Plan.DeviceLimit, &subscription.Plan.IsActive,
		&subscription.Plan.CreatedAt, &subscription.Plan.UpdatedAt)
	if err != nil {
		return nil, err
	}
	subscription.Plan.ID = subscription.PlanID
	return subscription, nil
}

// GetUserSubscriptions retrieves all subscriptions for a user
func (r *BillingRepository) GetUserSubscriptions(userID string) ([]*models.Subscription, error) {
	query := `
		SELECT s.id, s.user_id, s.plan_id, s.status, s.started_at, s.ends_at, s.created_at, s.updated_at,
		       p.name, p.description, p.price, p.billing_interval, p.device_limit, p.is_active, p.created_at, p.updated_at
		FROM subscriptions s
		JOIN pricing_plans p ON s.plan_id = p.id
		WHERE s.user_id = $1 ORDER BY s.created_at DESC
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptions []*models.Subscription
	for rows.Next() {
		subscription := &models.Subscription{Plan: &models.PricingPlan{}}
		err := rows.Scan(
			&subscription.ID, &subscription.UserID, &subscription.PlanID, &subscription.Status,
			&subscription.StartedAt, &subscription.EndsAt, &subscription.CreatedAt, &subscription.UpdatedAt,
			&subscription.Plan.Name, &subscription.Plan.Description, &subscription.Plan.Price,
			&subscription.Plan.BillingInterval, &subscription.Plan.DeviceLimit, &subscription.Plan.IsActive,
			&subscription.Plan.CreatedAt, &subscription.Plan.UpdatedAt)
		if err != nil {
			return nil, err
		}
		subscription.Plan.ID = subscription.PlanID
		subscriptions = append(subscriptions, subscription)
	}
	return subscriptions, nil
}

// UpdateSubscriptionStatus updates subscription status
func (r *BillingRepository) UpdateSubscriptionStatus(id, status string) error {
	query := `UPDATE subscriptions SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(query, status, time.Now(), id)
	return err
}

// CreatePayment creates a new payment record
func (r *BillingRepository) CreatePayment(payment *models.Payment) error {
	query := `
		INSERT INTO payments (id, subscription_id, amount, currency, status, payment_method, external_id, processed_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.Exec(query, payment.ID, payment.SubscriptionID, payment.Amount, payment.Currency,
		payment.Status, payment.PaymentMethod, payment.ExternalID, payment.ProcessedAt, payment.CreatedAt)
	return err
}

// UpdatePaymentStatus updates payment status
func (r *BillingRepository) UpdatePaymentStatus(id, status, externalID string) error {
	query := `UPDATE payments SET status = $1, external_id = $2, processed_at = $3 WHERE id = $4`
	processedAt := time.Now()
	_, err := r.db.Exec(query, status, externalID, processedAt, id)
	return err
}

// GetPaymentsBySubscription retrieves all payments for a subscription
func (r *BillingRepository) GetPaymentsBySubscription(subscriptionID string) ([]*models.Payment, error) {
	query := `
		SELECT id, subscription_id, amount, currency, status, payment_method, external_id, processed_at, created_at
		FROM payments WHERE subscription_id = $1 ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, subscriptionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []*models.Payment
	for rows.Next() {
		payment := &models.Payment{}
		err := rows.Scan(&payment.ID, &payment.SubscriptionID, &payment.Amount, &payment.Currency,
			&payment.Status, &payment.PaymentMethod, &payment.ExternalID, &payment.ProcessedAt, &payment.CreatedAt)
		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}
	return payments, nil
}

// GetBillingStats retrieves billing statistics
func (r *BillingRepository) GetBillingStats() (*models.BillingStats, error) {
	stats := &models.BillingStats{}
	
	// Active subscriptions count
	query := `SELECT COUNT(*) FROM subscriptions WHERE status = 'active'`
	err := r.db.QueryRow(query).Scan(&stats.TotalActiveSubscriptions)
	if err != nil {
		return nil, err
	}
	
	// Total revenue
	query = `SELECT COALESCE(SUM(amount), 0) FROM payments WHERE status = 'completed'`
	err = r.db.QueryRow(query).Scan(&stats.TotalRevenue)
	if err != nil {
		return nil, err
	}
	
	// Monthly revenue
	query = `SELECT COALESCE(SUM(amount), 0) FROM payments 
	         WHERE status = 'completed' AND processed_at >= DATE_TRUNC('month', CURRENT_DATE)`
	err = r.db.QueryRow(query).Scan(&stats.MonthlyRevenue)
	if err != nil {
		return nil, err
	}
	
	// Total payments
	query = `SELECT COUNT(*) FROM payments`
	err = r.db.QueryRow(query).Scan(&stats.TotalPayments)
	if err != nil {
		return nil, err
	}
	
	return stats, nil
}