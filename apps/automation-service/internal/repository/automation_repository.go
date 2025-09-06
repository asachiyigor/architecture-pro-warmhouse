package repository

import (
	"automation-service/internal/models"
	"database/sql"
	"time"
)

type AutomationRepository struct {
	db *sql.DB
}

func NewAutomationRepository(db *sql.DB) *AutomationRepository {
	return &AutomationRepository{db: db}
}

// CreateScenario creates a new automation scenario
func (r *AutomationRepository) CreateScenario(scenario *models.Scenario) error {
	query := `
		INSERT INTO scenarios (id, household_id, created_by, name, description, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.Exec(query, scenario.ID, scenario.HouseholdID, scenario.CreatedBy,
		scenario.Name, scenario.Description, scenario.IsActive, scenario.CreatedAt, scenario.UpdatedAt)
	return err
}

// GetScenarioByID retrieves a scenario by ID
func (r *AutomationRepository) GetScenarioByID(id string) (*models.Scenario, error) {
	query := `
		SELECT id, household_id, created_by, name, description, is_active, created_at, updated_at
		FROM scenarios WHERE id = $1
	`
	scenario := &models.Scenario{}
	err := r.db.QueryRow(query, id).Scan(
		&scenario.ID, &scenario.HouseholdID, &scenario.CreatedBy, &scenario.Name,
		&scenario.Description, &scenario.IsActive, &scenario.CreatedAt, &scenario.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return scenario, nil
}

// GetHouseholdScenarios retrieves all scenarios for a household
func (r *AutomationRepository) GetHouseholdScenarios(householdID string) ([]*models.Scenario, error) {
	query := `
		SELECT id, household_id, created_by, name, description, is_active, created_at, updated_at
		FROM scenarios WHERE household_id = $1 ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, householdID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scenarios []*models.Scenario
	for rows.Next() {
		scenario := &models.Scenario{}
		err := rows.Scan(&scenario.ID, &scenario.HouseholdID, &scenario.CreatedBy,
			&scenario.Name, &scenario.Description, &scenario.IsActive,
			&scenario.CreatedAt, &scenario.UpdatedAt)
		if err != nil {
			return nil, err
		}
		scenarios = append(scenarios, scenario)
	}
	return scenarios, nil
}

// UpdateScenarioStatus updates scenario active status
func (r *AutomationRepository) UpdateScenarioStatus(id string, isActive bool) error {
	query := `UPDATE scenarios SET is_active = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(query, isActive, time.Now(), id)
	return err
}

// CreateAction creates a new scenario action
func (r *AutomationRepository) CreateAction(action *models.ScenarioAction) error {
	query := `
		INSERT INTO scenario_actions (id, scenario_id, device_id, action_type, parameters, order_position, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(query, action.ID, action.ScenarioID, action.DeviceID,
		action.ActionType, action.Parameters, action.Order, action.CreatedAt)
	return err
}

// GetScenarioActions retrieves all actions for a scenario
func (r *AutomationRepository) GetScenarioActions(scenarioID string) ([]*models.ScenarioAction, error) {
	query := `
		SELECT id, scenario_id, device_id, action_type, parameters, order_position, created_at
		FROM scenario_actions WHERE scenario_id = $1 ORDER BY order_position ASC
	`
	rows, err := r.db.Query(query, scenarioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var actions []*models.ScenarioAction
	for rows.Next() {
		action := &models.ScenarioAction{}
		err := rows.Scan(&action.ID, &action.ScenarioID, &action.DeviceID,
			&action.ActionType, &action.Parameters, &action.Order, &action.CreatedAt)
		if err != nil {
			return nil, err
		}
		actions = append(actions, action)
	}
	return actions, nil
}

// CreateTrigger creates a new scenario trigger
func (r *AutomationRepository) CreateTrigger(trigger *models.ScenarioTrigger) error {
	query := `
		INSERT INTO scenario_triggers (id, scenario_id, trigger_type, device_id, condition, is_active, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(query, trigger.ID, trigger.ScenarioID, trigger.TriggerType,
		trigger.DeviceID, trigger.Condition, trigger.IsActive, trigger.CreatedAt)
	return err
}

// GetScenarioTriggers retrieves all triggers for a scenario
func (r *AutomationRepository) GetScenarioTriggers(scenarioID string) ([]*models.ScenarioTrigger, error) {
	query := `
		SELECT id, scenario_id, trigger_type, device_id, condition, is_active, created_at
		FROM scenario_triggers WHERE scenario_id = $1
	`
	rows, err := r.db.Query(query, scenarioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var triggers []*models.ScenarioTrigger
	for rows.Next() {
		trigger := &models.ScenarioTrigger{}
		err := rows.Scan(&trigger.ID, &trigger.ScenarioID, &trigger.TriggerType,
			&trigger.DeviceID, &trigger.Condition, &trigger.IsActive, &trigger.CreatedAt)
		if err != nil {
			return nil, err
		}
		triggers = append(triggers, trigger)
	}
	return triggers, nil
}

// CreateExecution creates a new execution record
func (r *AutomationRepository) CreateExecution(execution *models.ScenarioExecution) error {
	query := `
		INSERT INTO scenario_executions (id, scenario_id, trigger_type, status, started_at, completed_at, error_message)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(query, execution.ID, execution.ScenarioID, execution.TriggerType,
		execution.Status, execution.StartedAt, execution.CompletedAt, execution.ErrorMsg)
	return err
}

// UpdateExecutionStatus updates execution status
func (r *AutomationRepository) UpdateExecutionStatus(id, status string, errorMsg *string) error {
	completedAt := time.Now()
	query := `UPDATE scenario_executions SET status = $1, completed_at = $2, error_message = $3 WHERE id = $4`
	_, err := r.db.Exec(query, status, completedAt, errorMsg, id)
	return err
}

// GetExecutionHistory retrieves execution history for a scenario
func (r *AutomationRepository) GetExecutionHistory(scenarioID string, limit int) ([]*models.ScenarioExecution, error) {
	query := `
		SELECT id, scenario_id, trigger_type, status, started_at, completed_at, error_message
		FROM scenario_executions WHERE scenario_id = $1 ORDER BY started_at DESC LIMIT $2
	`
	rows, err := r.db.Query(query, scenarioID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var executions []*models.ScenarioExecution
	for rows.Next() {
		execution := &models.ScenarioExecution{}
		err := rows.Scan(&execution.ID, &execution.ScenarioID, &execution.TriggerType,
			&execution.Status, &execution.StartedAt, &execution.CompletedAt, &execution.ErrorMsg)
		if err != nil {
			return nil, err
		}
		executions = append(executions, execution)
	}
	return executions, nil
}

// GetAutomationStats retrieves automation statistics
func (r *AutomationRepository) GetAutomationStats() (*models.AutomationStats, error) {
	stats := &models.AutomationStats{}
	
	// Total scenarios
	query := `SELECT COUNT(*) FROM scenarios`
	err := r.db.QueryRow(query).Scan(&stats.TotalScenarios)
	if err != nil {
		return nil, err
	}
	
	// Active scenarios
	query = `SELECT COUNT(*) FROM scenarios WHERE is_active = true`
	err = r.db.QueryRow(query).Scan(&stats.ActiveScenarios)
	if err != nil {
		return nil, err
	}
	
	// Total executions
	query = `SELECT COUNT(*) FROM scenario_executions`
	err = r.db.QueryRow(query).Scan(&stats.TotalExecutions)
	if err != nil {
		return nil, err
	}
	
	// Successful executions
	query = `SELECT COUNT(*) FROM scenario_executions WHERE status = 'completed'`
	err = r.db.QueryRow(query).Scan(&stats.SuccessfulExecutions)
	if err != nil {
		return nil, err
	}
	
	// Failed executions
	query = `SELECT COUNT(*) FROM scenario_executions WHERE status = 'failed'`
	err = r.db.QueryRow(query).Scan(&stats.FailedExecutions)
	if err != nil {
		return nil, err
	}
	
	// Executions today
	query = `SELECT COUNT(*) FROM scenario_executions WHERE DATE(started_at) = CURRENT_DATE`
	err = r.db.QueryRow(query).Scan(&stats.ExecutionsToday)
	if err != nil {
		return nil, err
	}
	
	return stats, nil
}