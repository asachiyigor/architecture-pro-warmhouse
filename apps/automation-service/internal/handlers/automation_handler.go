package handlers

import (
	"automation-service/internal/models"
	"automation-service/internal/service"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AutomationHandler struct {
	automationService *service.AutomationService
}

func NewAutomationHandler(automationService *service.AutomationService) *AutomationHandler {
	return &AutomationHandler{automationService: automationService}
}

// RegisterRoutes registers all automation-related routes
func (h *AutomationHandler) RegisterRoutes(router *gin.RouterGroup) {
	scenarios := router.Group("/scenarios")
	{
		scenarios.POST("", h.CreateScenario)
		scenarios.GET("/:id", h.GetScenario)
		scenarios.PUT("/:id/activate", h.ActivateScenario)
		scenarios.PUT("/:id/deactivate", h.DeactivateScenario)
		scenarios.GET("/household/:household_id", h.GetHouseholdScenarios)
		
		scenarios.POST("/:id/actions", h.AddAction)
		scenarios.POST("/:id/triggers", h.AddTrigger)
		scenarios.POST("/:id/execute", h.ExecuteScenario)
		scenarios.GET("/:id/executions", h.GetExecutionHistory)
	}

	router.GET("/stats", h.GetAutomationStats)
}

// CreateScenario creates a new automation scenario
// @Summary Create New Automation Scenario
// @Description Создать новый сценарий автоматизации
// @Tags Scenarios
// @Accept json
// @Produce json
// @Param scenario body models.CreateScenarioRequest true "Scenario data"
// @Success 201 {object} models.Scenario
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/scenarios [post]
func (h *AutomationHandler) CreateScenario(c *gin.Context) {
	var req models.CreateScenarioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	scenario, err := h.automationService.CreateScenario(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, scenario)
}

// GetScenario retrieves a scenario by ID
// @Summary Get Scenario by ID
// @Description Получить сценарий по ID
// @Tags Scenarios
// @Produce json
// @Param id path string true "Scenario ID"
// @Success 200 {object} models.Scenario
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/scenarios/{id} [get]
func (h *AutomationHandler) GetScenario(c *gin.Context) {
	id := c.Param("id")

	scenario, err := h.automationService.GetScenario(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Scenario not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, scenario)
}

// GetHouseholdScenarios retrieves all scenarios for a household
// @Summary Get Household Scenarios
// @Description Получить все сценарии для домохозяйства
// @Tags Scenarios
// @Produce json
// @Param household_id path string true "Household ID"
// @Success 200 {array} models.Scenario
// @Failure 500 {object} map[string]string
// @Router /api/v1/scenarios/household/{household_id} [get]
func (h *AutomationHandler) GetHouseholdScenarios(c *gin.Context) {
	householdID := c.Param("household_id")

	scenarios, err := h.automationService.GetHouseholdScenarios(householdID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, scenarios)
}

// ActivateScenario activates a scenario
// @Summary Activate Scenario
// @Description Активировать сценарий автоматизации
// @Tags Scenarios
// @Produce json
// @Param id path string true "Scenario ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/scenarios/{id}/activate [put]
func (h *AutomationHandler) ActivateScenario(c *gin.Context) {
	id := c.Param("id")

	err := h.automationService.ActivateScenario(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "scenario not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Scenario activated successfully"})
}

// DeactivateScenario deactivates a scenario
func (h *AutomationHandler) DeactivateScenario(c *gin.Context) {
	id := c.Param("id")

	err := h.automationService.DeactivateScenario(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "scenario not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Scenario deactivated successfully"})
}

// AddAction adds an action to a scenario
func (h *AutomationHandler) AddAction(c *gin.Context) {
	scenarioID := c.Param("id")
	
	var req models.CreateActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	action, err := h.automationService.AddAction(scenarioID, req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "scenario not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, action)
}

// AddTrigger adds a trigger to a scenario
func (h *AutomationHandler) AddTrigger(c *gin.Context) {
	scenarioID := c.Param("id")
	
	var req models.CreateTriggerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	trigger, err := h.automationService.AddTrigger(scenarioID, req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "scenario not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, trigger)
}

// ExecuteScenario manually executes a scenario
// @Summary Execute Automation Scenario
// @Description Выполнить сценарий автоматизации
// @Tags Scenarios
// @Accept json
// @Produce json
// @Param id path string true "Scenario ID"
// @Param execution body object{user_id=string} true "Execution data"
// @Success 202 {object} models.ScenarioExecution
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/scenarios/{id}/execute [post]
func (h *AutomationHandler) ExecuteScenario(c *gin.Context) {
	scenarioID := c.Param("id")
	
	var requestBody struct {
		UserID string `json:"user_id" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req := models.ExecuteScenarioRequest{
		ScenarioID: scenarioID,
		UserID:     requestBody.UserID,
	}

	execution, err := h.automationService.ExecuteScenario(req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "scenario not found" {
			status = http.StatusNotFound
		} else if err.Error() == "scenario is not active" {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, execution)
}

// GetExecutionHistory retrieves execution history for a scenario
func (h *AutomationHandler) GetExecutionHistory(c *gin.Context) {
	scenarioID := c.Param("id")
	
	// Parse limit from query parameter, default to 10
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	executions, err := h.automationService.GetExecutionHistory(scenarioID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, executions)
}

// GetAutomationStats retrieves automation statistics
// @Summary Get Automation Statistics
// @Description Получить статистику автоматизации
// @Tags Statistics
// @Produce json
// @Success 200 {object} models.AutomationStats
// @Failure 500 {object} map[string]string
// @Router /api/v1/stats [get]
func (h *AutomationHandler) GetAutomationStats(c *gin.Context) {
	stats, err := h.automationService.GetAutomationStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}