package handlers

import (
	"billing-service/internal/models"
	"billing-service/internal/service"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BillingHandler struct {
	billingService *service.BillingService
}

func NewBillingHandler(billingService *service.BillingService) *BillingHandler {
	return &BillingHandler{billingService: billingService}
}

// RegisterRoutes registers all billing-related routes
func (h *BillingHandler) RegisterRoutes(router *gin.RouterGroup) {
	plans := router.Group("/plans")
	{
		plans.GET("", h.GetPlans)
		plans.GET("/:id", h.GetPlan)
	}

	subscriptions := router.Group("/subscriptions")
	{
		subscriptions.POST("", h.CreateSubscription)
		subscriptions.GET("/:id", h.GetSubscription)
		subscriptions.PUT("/:id/cancel", h.CancelSubscription)
		subscriptions.GET("/user/:user_id", h.GetUserSubscriptions)
	}

	payments := router.Group("/payments")
	{
		payments.POST("", h.ProcessPayment)
		payments.GET("/subscription/:subscription_id", h.GetPaymentHistory)
	}

	router.GET("/stats", h.GetBillingStats)
}

// GetPlans retrieves all pricing plans
// @Summary Get All Pricing Plans
// @Description Получить все тарифные планы
// @Tags Plans
// @Produce json
// @Success 200 {array} models.PricingPlan
// @Failure 500 {object} map[string]string
// @Router /api/v1/plans [get]
func (h *BillingHandler) GetPlans(c *gin.Context) {
	plans, err := h.billingService.GetPlans()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, plans)
}

// GetPlan retrieves a specific pricing plan
// @Summary Get Pricing Plan by ID
// @Description Получить тарифный план по ID
// @Tags Plans
// @Produce json
// @Param id path string true "Plan ID"
// @Success 200 {object} models.PricingPlan
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/plans/{id} [get]
func (h *BillingHandler) GetPlan(c *gin.Context) {
	id := c.Param("id")

	plan, err := h.billingService.GetPlan(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Plan not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, plan)
}

// CreateSubscription creates a new subscription
// @Summary Create New Subscription
// @Description Создать новую подписку
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Param subscription body models.CreateSubscriptionRequest true "Subscription data"
// @Success 201 {object} models.Subscription
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/subscriptions [post]
func (h *BillingHandler) CreateSubscription(c *gin.Context) {
	var req models.CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subscription, err := h.billingService.CreateSubscription(req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "pricing plan not found" {
			status = http.StatusNotFound
		} else if err.Error() == "pricing plan is not active" || err.Error() == "user already has an active subscription" {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, subscription)
}

// GetSubscription retrieves a subscription by ID
func (h *BillingHandler) GetSubscription(c *gin.Context) {
	id := c.Param("id")

	subscription, err := h.billingService.GetSubscription(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subscription)
}

// GetUserSubscriptions retrieves all subscriptions for a user
func (h *BillingHandler) GetUserSubscriptions(c *gin.Context) {
	userID := c.Param("user_id")

	subscriptions, err := h.billingService.GetUserSubscriptions(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subscriptions)
}

// CancelSubscription cancels a subscription
func (h *BillingHandler) CancelSubscription(c *gin.Context) {
	id := c.Param("id")

	err := h.billingService.CancelSubscription(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err == sql.ErrNoRows {
			status = http.StatusNotFound
		} else if err.Error() == "subscription is not active" {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscription cancelled successfully"})
}

// ProcessPayment processes a payment
// @Summary Process Payment
// @Description Обработать платеж
// @Tags Payments
// @Accept json
// @Produce json
// @Param payment body models.PaymentRequest true "Payment data"
// @Success 201 {object} models.Payment
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/payments [post]
func (h *BillingHandler) ProcessPayment(c *gin.Context) {
	var req models.PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, err := h.billingService.ProcessPayment(req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "subscription not found" {
			status = http.StatusNotFound
		} else if err.Error() == "subscription is not active" {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, payment)
}

// GetPaymentHistory retrieves payment history for a subscription
func (h *BillingHandler) GetPaymentHistory(c *gin.Context) {
	subscriptionID := c.Param("subscription_id")

	payments, err := h.billingService.GetPaymentHistory(subscriptionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payments)
}

// GetBillingStats retrieves billing statistics
// @Summary Get Billing Statistics
// @Description Получить статистику платежей и подписок
// @Tags Statistics
// @Produce json
// @Success 200 {object} models.BillingStats
// @Failure 500 {object} map[string]string
// @Router /api/v1/stats [get]
func (h *BillingHandler) GetBillingStats(c *gin.Context) {
	stats, err := h.billingService.GetBillingStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}