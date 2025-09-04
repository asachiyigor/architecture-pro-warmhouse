package models

import (
	"time"
	"github.com/google/uuid"
)

// Scenario represents an automation scenario
type Scenario struct {
	ID          string    `json:"id" db:"id"`
	HouseholdID string    `json:"household_id" db:"household_id"`
	CreatedBy   string    `json:"created_by" db:"created_by"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`

	// Embedded actions for convenience
	Actions []*ScenarioAction `json:"actions,omitempty"`
	Triggers []*ScenarioTrigger `json:"triggers,omitempty"`
}

// ScenarioAction represents an action within a scenario
type ScenarioAction struct {
	ID          string            `json:"id" db:"id"`
	ScenarioID  string            `json:"scenario_id" db:"scenario_id"`
	DeviceID    string            `json:"device_id" db:"device_id"`
	ActionType  string            `json:"action_type" db:"action_type"`
	Parameters  string            `json:"parameters" db:"parameters"` // JSON string
	Order       int               `json:"order" db:"order_position"`
	CreatedAt   time.Time         `json:"created_at" db:"created_at"`
}

// ScenarioTrigger represents a trigger for a scenario
type ScenarioTrigger struct {
	ID           string    `json:"id" db:"id"`
	ScenarioID   string    `json:"scenario_id" db:"scenario_id"`
	TriggerType  string    `json:"trigger_type" db:"trigger_type"` // time, device_state, manual
	DeviceID     *string   `json:"device_id" db:"device_id"`
	Condition    string    `json:"condition" db:"condition"`    // JSON string with condition details
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// ScenarioExecution represents execution history
type ScenarioExecution struct {
	ID          string     `json:"id" db:"id"`
	ScenarioID  string     `json:"scenario_id" db:"scenario_id"`
	TriggerType string     `json:"trigger_type" db:"trigger_type"`
	Status      string     `json:"status" db:"status"` // started, completed, failed, cancelled
	StartedAt   time.Time  `json:"started_at" db:"started_at"`
	CompletedAt *time.Time `json:"completed_at" db:"completed_at"`
	ErrorMsg    *string    `json:"error_message" db:"error_message"`
}

// CreateScenarioRequest represents request to create a scenario
type CreateScenarioRequest struct {
	HouseholdID string `json:"household_id" binding:"required"`
	CreatedBy   string `json:"created_by" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// CreateActionRequest represents request to create an action
type CreateActionRequest struct {
	DeviceID   string                 `json:"device_id" binding:"required"`
	ActionType string                 `json:"action_type" binding:"required"`
	Parameters map[string]interface{} `json:"parameters"`
	Order      int                    `json:"order"`
}

// CreateTriggerRequest represents request to create a trigger
type CreateTriggerRequest struct {
	TriggerType string                 `json:"trigger_type" binding:"required"`
	DeviceID    *string                `json:"device_id"`
	Condition   map[string]interface{} `json:"condition" binding:"required"`
}

// ExecuteScenarioRequest represents manual execution request
type ExecuteScenarioRequest struct {
	ScenarioID string `json:"scenario_id" binding:"required"`
	UserID     string `json:"user_id" binding:"required"`
}

// AutomationStats represents automation statistics
type AutomationStats struct {
	TotalScenarios       int `json:"total_scenarios"`
	ActiveScenarios      int `json:"active_scenarios"`
	TotalExecutions      int `json:"total_executions"`
	SuccessfulExecutions int `json:"successful_executions"`
	FailedExecutions     int `json:"failed_executions"`
	ExecutionsToday      int `json:"executions_today"`
}

// NewScenario creates a new scenario with generated ID
func NewScenario(req CreateScenarioRequest) *Scenario {
	now := time.Now()
	return &Scenario{
		ID:          uuid.New().String(),
		HouseholdID: req.HouseholdID,
		CreatedBy:   req.CreatedBy,
		Name:        req.Name,
		Description: req.Description,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// NewScenarioAction creates a new scenario action
func NewScenarioAction(scenarioID string, req CreateActionRequest) *ScenarioAction {
	parametersJSON := "{}"
	if req.Parameters != nil {
		// In real implementation, properly marshal JSON
		parametersJSON = `{"parameters": "example"}`
	}
	
	return &ScenarioAction{
		ID:         uuid.New().String(),
		ScenarioID: scenarioID,
		DeviceID:   req.DeviceID,
		ActionType: req.ActionType,
		Parameters: parametersJSON,
		Order:      req.Order,
		CreatedAt:  time.Now(),
	}
}

// NewScenarioTrigger creates a new scenario trigger
func NewScenarioTrigger(scenarioID string, req CreateTriggerRequest) *ScenarioTrigger {
	conditionJSON := `{"condition": "example"}`
	// In real implementation, properly marshal JSON
	
	return &ScenarioTrigger{
		ID:          uuid.New().String(),
		ScenarioID:  scenarioID,
		TriggerType: req.TriggerType,
		DeviceID:    req.DeviceID,
		Condition:   conditionJSON,
		IsActive:    true,
		CreatedAt:   time.Now(),
	}
}

// NewScenarioExecution creates a new execution record
func NewScenarioExecution(scenarioID, triggerType string) *ScenarioExecution {
	return &ScenarioExecution{
		ID:          uuid.New().String(),
		ScenarioID:  scenarioID,
		TriggerType: triggerType,
		Status:      "started",
		StartedAt:   time.Now(),
	}
}