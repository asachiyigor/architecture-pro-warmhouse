package service

import (
	"database/sql"
	"errors"
	"user-service/internal/models"
	"user-service/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(req models.CreateUserRequest) (*models.User, error) {
	// Check if user with this email already exists
	existingUser, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	user := models.NewUser(req)
	err = s.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(id string) (*models.User, error) {
	return s.userRepo.GetUserByID(id)
}

// ValidateUser checks if user exists and is active
func (s *UserService) ValidateUser(id string) (*models.User, error) {
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	if !user.IsActive {
		return nil, errors.New("user is inactive")
	}
	return user, nil
}

// CreateHousehold creates a new household for a user
func (s *UserService) CreateHousehold(ownerID string, req models.CreateHouseholdRequest) (*models.Household, error) {
	// Validate that owner exists
	_, err := s.ValidateUser(ownerID)
	if err != nil {
		return nil, err
	}

	household := models.NewHousehold(ownerID, req)
	err = s.userRepo.CreateHousehold(household)
	if err != nil {
		return nil, err
	}

	return household, nil
}

// GetHousehold retrieves a household by ID
func (s *UserService) GetHousehold(id string) (*models.Household, error) {
	return s.userRepo.GetHouseholdByID(id)
}

// GetUserHouseholds retrieves all households owned by a user
func (s *UserService) GetUserHouseholds(ownerID string) ([]*models.Household, error) {
	return s.userRepo.GetHouseholdsByOwnerID(ownerID)
}

// CheckHouseholdAccess verifies that user has access to household
func (s *UserService) CheckHouseholdAccess(userID, householdID string) error {
	return s.userRepo.CheckHouseholdAccess(userID, householdID)
}