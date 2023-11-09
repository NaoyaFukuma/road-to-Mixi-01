// repository/user.go

package repository

import (
	"database/sql"
	"minimal_sns_app/domain/models"
	"minimal_sns_app/logutils"
)

// UserRepository defines the interface for user data access.
type UserRepository interface {
	GetUser(userID int64) (*models.User, error)
	CreateUser(name string) (*models.User, error)
	DeleteUser(userID int64) error
}

type userRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new instance of a UserRepository.
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

// GetUser retrieves a user by ID.
func (r *userRepository) GetUser(userID int64) (*models.User, error) {
	var user models.User
	query := `SELECT id, name FROM users WHERE id = ?`

	err := r.db.QueryRow(query, userID).Scan(&user.ID, &user.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logutils.Error(err.Error())
		return nil, err
	}
	return &user, nil
}

// CreateUser creates a new user.
func (r *userRepository) CreateUser(name string) (*models.User, error) {
	query := `INSERT INTO users (name) VALUES (?)`
	result, err := r.db.Exec(query, name)
	if err != nil {
		logutils.Error(err.Error())
		return nil, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		logutils.Error(err.Error())
		return nil, err
	}

	return r.GetUser(userID)
}

// DeleteUser deletes a user by ID.
func (r *userRepository) DeleteUser(userID int64) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := r.db.Exec(query, userID)
	if err != nil {
		logutils.Error(err.Error())
		return err
	}
	return nil
}
