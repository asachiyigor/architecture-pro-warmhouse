package repository

import (
	"database/sql"
	"user-service/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser creates a new user in the database
func (r *UserRepository) CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (id, email, first_name, last_name, phone, timezone, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.Exec(query, user.ID, user.Email, user.FirstName, user.LastName, 
		user.Phone, user.Timezone, user.IsActive, user.CreatedAt, user.UpdatedAt)
	return err
}

// GetUserByID retrieves a user by ID
func (r *UserRepository) GetUserByID(id string) (*models.User, error) {
	query := `
		SELECT id, email, first_name, last_name, phone, timezone, is_active, created_at, updated_at
		FROM users WHERE id = $1
	`
	user := &models.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.FirstName, &user.LastName,
		&user.Phone, &user.Timezone, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByEmail retrieves a user by email
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, email, first_name, last_name, phone, timezone, is_active, created_at, updated_at
		FROM users WHERE email = $1
	`
	user := &models.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.FirstName, &user.LastName,
		&user.Phone, &user.Timezone, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// CreateHousehold creates a new household
func (r *UserRepository) CreateHousehold(household *models.Household) error {
	query := `
		INSERT INTO households (id, owner_id, name, address, timezone, coordinates, wifi_ssid, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.Exec(query, household.ID, household.OwnerID, household.Name, household.Address,
		household.Timezone, household.Coordinates, household.WiFiSSID, household.CreatedAt, household.UpdatedAt)
	return err
}

// GetHouseholdByID retrieves a household by ID
func (r *UserRepository) GetHouseholdByID(id string) (*models.Household, error) {
	query := `
		SELECT id, owner_id, name, address, timezone, coordinates, wifi_ssid, created_at, updated_at
		FROM households WHERE id = $1
	`
	household := &models.Household{}
	err := r.db.QueryRow(query, id).Scan(
		&household.ID, &household.OwnerID, &household.Name, &household.Address,
		&household.Timezone, &household.Coordinates, &household.WiFiSSID,
		&household.CreatedAt, &household.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return household, nil
}

// GetHouseholdsByOwnerID retrieves all households owned by a user
func (r *UserRepository) GetHouseholdsByOwnerID(ownerID string) ([]*models.Household, error) {
	query := `
		SELECT id, owner_id, name, address, timezone, coordinates, wifi_ssid, created_at, updated_at
		FROM households WHERE owner_id = $1
	`
	rows, err := r.db.Query(query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var households []*models.Household
	for rows.Next() {
		household := &models.Household{}
		err := rows.Scan(
			&household.ID, &household.OwnerID, &household.Name, &household.Address,
			&household.Timezone, &household.Coordinates, &household.WiFiSSID,
			&household.CreatedAt, &household.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		households = append(households, household)
	}
	return households, nil
}

// CheckHouseholdAccess checks if user has access to household
func (r *UserRepository) CheckHouseholdAccess(userID, householdID string) error {
	query := `
		SELECT 1 FROM households h
		LEFT JOIN household_members hm ON h.id = hm.household_id
		WHERE h.id = $1 AND (h.owner_id = $2 OR hm.user_id = $2)
		LIMIT 1
	`
	var exists int
	err := r.db.QueryRow(query, householdID, userID).Scan(&exists)
	return err
}