package services

import (
	"fmt"
	"time"

	"github.com/CEM-KEA/whoknows/backend/internal/models"
	"github.com/CEM-KEA/whoknows/backend/internal/security"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// CreateUser creates a new user in the database.
func CreateUser(db *gorm.DB, user *models.User) error {
	return db.Create(user).Error
}

// UpdateUser updates a user in the database.
func UpdateUser(db *gorm.DB, user *models.User) error {
	return db.Save(user).Error
}

// GetUserByUsername retrieves a user from the database by username.
func GetUserByUsername(db *gorm.DB, username string) (*models.User, error) {
	user := &models.User{}
	err := db.Where("username = ?", username).First(user).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user by username")
	}
	return user, nil
}

// GetUserByID retrieves a user by their ID.
func GetUserByID(db *gorm.DB, userID uint) (*models.User, error) {
	user := &models.User{}
	if err := db.First(user, userID).Error; err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// UpdateLastLogin updates the last_login field for a user.
func UpdateLastLogin(db *gorm.DB, user *models.User) error {
	user.LastLogin = time.Now()
	return db.Save(user).Error
}

// CheckUserPassword checks if the provided password matches the user's password hash.
func CheckUserPassword(db *gorm.DB, password string, username string) (*models.User, bool, error) {
	user, err := GetUserByUsername(db, username)
	if err != nil {
		return user, false, fmt.Errorf("invalid username")
	}
	if !security.CheckPasswordHash(password, user.PasswordHash) {
		return user, false, fmt.Errorf("invalid password")
	}
	return user, true, nil
}
