package handlers

import (
	"database/sql"
	"net/http"
	"user-service/internal/models"
	"user-service/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// RegisterRoutes registers all user-related routes
func (h *UserHandler) RegisterRoutes(router *gin.RouterGroup) {
	users := router.Group("/users")
	{
		users.POST("", h.CreateUser)
		users.GET("/:id", h.GetUser)
		users.POST("/:id/validate", h.ValidateUser)
	}

	households := router.Group("/households")
	{
		households.POST("", h.CreateHousehold)
		households.GET("/:id", h.GetHousehold)
		households.GET("/owner/:owner_id", h.GetUserHouseholds)
		households.POST("/:id/access/:user_id", h.CheckHouseholdAccess)
	}
}

// CreateUser creates a new user
// @Summary Создать пользователя
// @Description Создает нового пользователя в системе
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.CreateUserRequest true "Данные пользователя"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.CreateUser(req)
	if err != nil {
		if err.Error() == "user with this email already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// GetUser retrieves a user by ID
// @Summary Получить пользователя по ID
// @Description Возвращает данные пользователя по его ID
// @Tags Users
// @Produce json
// @Param id path string true "ID пользователя"
// @Success 200 {object} models.User
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")

	user, err := h.userService.GetUser(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// ValidateUser validates a user by ID
func (h *UserHandler) ValidateUser(c *gin.Context) {
	id := c.Param("id")

	user, err := h.userService.ValidateUser(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		if err.Error() == "user is inactive" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid": true,
		"user":  user,
	})
}

// CreateHousehold creates a new household
// @Summary Создать домохозяйство
// @Description Создает новое домохозяйство для пользователя
// @Tags Households
// @Accept json
// @Produce json
// @Param household body object true "Данные домохозяйства"
// @Success 201 {object} models.Household
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /households [post]
func (h *UserHandler) CreateHousehold(c *gin.Context) {
	var req struct {
		OwnerID string                      `json:"owner_id" binding:"required"`
		models.CreateHouseholdRequest
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	household, err := h.userService.CreateHousehold(req.OwnerID, req.CreateHouseholdRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, household)
}

// GetHousehold retrieves a household by ID
// @Summary Получить домохозяйство по ID
// @Description Возвращает данные домохозяйства по его ID
// @Tags Households
// @Produce json
// @Param id path string true "ID домохозяйства"
// @Success 200 {object} models.Household
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /households/{id} [get]
func (h *UserHandler) GetHousehold(c *gin.Context) {
	id := c.Param("id")

	household, err := h.userService.GetHousehold(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Household not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, household)
}

// GetUserHouseholds retrieves all households for a user
func (h *UserHandler) GetUserHouseholds(c *gin.Context) {
	ownerID := c.Param("owner_id")

	households, err := h.userService.GetUserHouseholds(ownerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, households)
}

// CheckHouseholdAccess checks if user has access to household
func (h *UserHandler) CheckHouseholdAccess(c *gin.Context) {
	householdID := c.Param("id")
	userID := c.Param("user_id")

	err := h.userService.CheckHouseholdAccess(userID, householdID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access": "granted"})
}