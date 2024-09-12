package services

import (
	"github.com/CEM-KEA/whoknows/backend/internal/models"
	"github.com/CEM-KEA/whoknows/backend/internal/database"

	"github.com/pkg/errors"
)

// CreateUser creates a new user in the database.
func CreateUser(user *models.User) error {
	return database.DB.Create(user).Error
}

// GetUserByEmail retrieves a user from the database by email.
func GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := database.DB.Where("email = ?", email).First(user).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user by email")
	}
	return user, nil
}