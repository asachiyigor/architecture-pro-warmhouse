package service

import (
	"automation-service/internal/models"
	"automation-service/internal/repository"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

type AutomationService struct {
	automationRepo *repository.AutomationRepository
}

func NewAutomationService(automationRepo *repository.AutomationRepository) *AutomationService {
	return &AutomationService{automationRepo: automationRepo}
}

// CreateScenario creates a new automation scenario
func (s *AutomationService) CreateScenario(req models.CreateScenarioRequest) (*models.Scenario, error) {
	scenario := models.NewScenario(req)
	err := s.automationRepo.CreateScenario(scenario)
	if err != nil {
		return nil, err
	}
	return scenario, nil
}

// GetScenario retrieves a scenario by ID with its actions and triggers
func (s *AutomationService) GetScenario(id string) (*models.Scenario, error) {
	scenario, err := s.automationRepo.GetScenarioByID(id)
	if err != nil {
		return nil, err
	}

	// Load actions
	actions, err := s.automationRepo.GetScenarioActions(id)
	if err != nil {
		return nil, err
	}
	scenario.Actions = actions

	// Load triggers
	triggers, err := s.automationRepo.GetScenarioTriggers(id)
	if err != nil {
		return nil, err
	}
	scenario.Triggers = triggers

	return scenario, nil
}

// GetHouseholdScenarios retrieves all scenarios for a household
func (s *AutomationService) GetHouseholdScenarios(householdID string) ([]*models.Scenario, error) {
	return s.automationRepo.GetHouseholdScenarios(householdID)
}

// ActivateScenario activates a scenario
func (s *AutomationService) ActivateScenario(id string) error {
	// Verify scenario exists
	_, err := s.automationRepo.GetScenarioByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("scenario not found")
		}
		return err
	}

	return s.automationRepo.UpdateScenarioStatus(id, true)
}

// DeactivateScenario deactivates a scenario
func (s *AutomationService) DeactivateScenario(id string) error {
	// Verify scenario exists
	_, err := s.automationRepo.GetScenarioByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("scenario not found")
		}
		return err
	}

	return s.automationRepo.UpdateScenarioStatus(id, false)
}

// AddAction adds an action to a scenario
func (s *AutomationService) AddAction(scenarioID string, req models.CreateActionRequest) (*models.ScenarioAction, error) {
	// Verify scenario exists
	_, err := s.automationRepo.GetScenarioByID(scenarioID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("scenario not found")
		}
		return nil, err
	}

	action := models.NewScenarioAction(scenarioID, req)
	err = s.automationRepo.CreateAction(action)
	if err != nil {
		return nil, err
	}

	return action, nil
}

// AddTrigger adds a trigger to a scenario
func (s *AutomationService) AddTrigger(scenarioID string, req models.CreateTriggerRequest) (*models.ScenarioTrigger, error) {
	// Verify scenario exists
	_, err := s.automationRepo.GetScenarioByID(scenarioID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("scenario not found")
		}
		return nil, err
	}

	trigger := models.NewScenarioTrigger(scenarioID, req)
	err = s.automationRepo.CreateTrigger(trigger)
	if err != nil {
		return nil, err
	}

	return trigger, nil
}

// ExecuteScenario manually executes a scenario
func (s *AutomationService) ExecuteScenario(req models.ExecuteScenarioRequest) (*models.ScenarioExecution, error) {
	// Verify scenario exists and is active
	scenario, err := s.automationRepo.GetScenarioByID(req.ScenarioID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("scenario not found")
		}
		return nil, err
	}

	if !scenario.IsActive {
		return nil, errors.New("scenario is not active")
	}

	// Create execution record
	execution := models.NewScenarioExecution(req.ScenarioID, "manual")
	err = s.automationRepo.CreateExecution(execution)
	if err != nil {
		return nil, err
	}

	// Execute scenario actions asynchronously
	go s.executeScenarioActions(execution)

	return execution, nil
}

// GetExecutionHistory retrieves execution history for a scenario
func (s *AutomationService) GetExecutionHistory(scenarioID string, limit int) ([]*models.ScenarioExecution, error) {
	return s.automationRepo.GetExecutionHistory(scenarioID, limit)
}

// GetAutomationStats retrieves automation statistics
func (s *AutomationService) GetAutomationStats() (*models.AutomationStats, error) {
	return s.automationRepo.GetAutomationStats()
}

// executeScenarioActions executes all actions in a scenario
func (s *AutomationService) executeScenarioActions(execution *models.ScenarioExecution) {
	log.Printf("Starting execution of scenario %s", execution.ScenarioID)

	// Get scenario actions
	actions, err := s.automationRepo.GetScenarioActions(execution.ScenarioID)
	if err != nil {
		errorMsg := fmt.Sprintf("Failed to get scenario actions: %v", err)
		s.automationRepo.UpdateExecutionStatus(execution.ID, "failed", &errorMsg)
		return
	}

	// Execute each action in order
	for _, action := range actions {
		err := s.executeAction(action)
		if err != nil {
			errorMsg := fmt.Sprintf("Failed to execute action %s: %v", action.ID, err)
			s.automationRepo.UpdateExecutionStatus(execution.ID, "failed", &errorMsg)
			return
		}

		// Small delay between actions
		time.Sleep(500 * time.Millisecond)
	}

	// Mark execution as completed
	s.automationRepo.UpdateExecutionStatus(execution.ID, "completed", nil)
	log.Printf("Completed execution of scenario %s", execution.ScenarioID)
}

// executeAction executes a single action
func (s *AutomationService) executeAction(action *models.ScenarioAction) error {
	log.Printf("Executing action %s: %s on device %s", action.ID, action.ActionType, action.DeviceID)

	// Simulate action execution
	// In real implementation, this would call the device service
	switch action.ActionType {
	case "turn_on":
		log.Printf("Turning on device %s", action.DeviceID)
	case "turn_off":
		log.Printf("Turning off device %s", action.DeviceID)
	case "set_temperature":
		log.Printf("Setting temperature on device %s", action.DeviceID)
	case "set_brightness":
		log.Printf("Setting brightness on device %s", action.DeviceID)
	default:
		return fmt.Errorf("unknown action type: %s", action.ActionType)
	}

	// Simulate some execution time
	time.Sleep(200 * time.Millisecond)
	return nil
}