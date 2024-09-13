package services

import (
	"github.com/CEM-KEA/whoknows/backend/internal/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// CreateUser creates a new user in the database.
func CreateUser(db *gorm.DB, user *models.User) error {
	return db.Create(user).Error
}

// GetUserByEmail retrieves a user from the database by email.
func GetUserByEmail(db *gorm.DB, email string) (*models.User, error) {
	user := &models.User{}
	err := db.Where("email = ?", email).First(user).Error

	if err != nil {
		return nil, errors.Wrap(err, "failed to get user by email")
	}

	return user, nil
}
