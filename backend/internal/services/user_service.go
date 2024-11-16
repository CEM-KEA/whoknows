package services

import (
	"time"

	"github.com/CEM-KEA/whoknows/backend/internal/models"
	"github.com/CEM-KEA/whoknows/backend/internal/security"
	"github.com/CEM-KEA/whoknows/backend/internal/utils"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// CreateUser creates a new user in the database.
func CreateUser(db *gorm.DB, user *models.User) error {
	utils.LogInfo("Creating new user", logrus.Fields{
		"username": user.Username,
		"email":    user.Email,
	})

	err := db.Create(user).Error
	if err != nil {
		utils.LogError(err, "Failed to create user", logrus.Fields{
			"username": user.Username,
			"email":    user.Email,
		})
		return err
	}

	utils.LogInfo("User created successfully", logrus.Fields{
		"username": user.Username,
		"email":    user.Email,
	})
	return nil
}

// UpdateUser updates a user in the database.
func UpdateUser(db *gorm.DB, user *models.User) error {
	utils.LogInfo("Updating user", logrus.Fields{
		"username": user.Username,
		"id":       user.ID,
	})

	err := db.Save(user).Error
	if err != nil {
		utils.LogError(err, "Failed to update user", logrus.Fields{
			"username": user.Username,
			"id":       user.ID,
		})
		return err
	}

	utils.LogInfo("User updated successfully", logrus.Fields{
		"username": user.Username,
		"id":       user.ID,
	})
	return nil
}

// GetUserByUsername retrieves a user from the database by username.
func GetUserByUsername(db *gorm.DB, username string) (*models.User, error) {
	utils.LogInfo("Retrieving user by username", logrus.Fields{
		"username": username,
	})

	user := &models.User{}
	err := db.Where("username = ?", username).First(user).Error
	if err != nil {
		utils.LogError(err, "Failed to retrieve user by username", logrus.Fields{
			"username": username,
		})
		return nil, errors.Wrap(err, "failed to get user by username")
	}

	utils.LogInfo("User retrieved successfully", logrus.Fields{
		"username": user.Username,
		"id":       user.ID,
	})
	return user, nil
}

// GetUserByID retrieves a user by their ID.
func GetUserByID(db *gorm.DB, userID uint) (*models.User, error) {
	utils.LogInfo("Retrieving user by ID", logrus.Fields{
		"userID": userID,
	})

	user := &models.User{}
	err := db.First(user, userID).Error
	if err != nil {
		utils.LogError(err, "Failed to retrieve user by ID", logrus.Fields{
			"userID": userID,
		})
		return nil, errors.New("user not found")
	}

	utils.LogInfo("User retrieved successfully", logrus.Fields{
		"userID": userID,
		"username": user.Username,
	})
	return user, nil
}

// UpdateLastLogin updates the last_login field for a user.
func UpdateLastLogin(db *gorm.DB, user *models.User) error {
	utils.LogInfo("Updating last login for user", logrus.Fields{
		"username": user.Username,
		"id":       user.ID,
	})

	user.LastLogin = time.Now()
	err := db.Save(user).Error
	if err != nil {
		utils.LogError(err, "Failed to update last login", logrus.Fields{
			"username": user.Username,
			"id":       user.ID,
		})
		return err
	}

	utils.LogInfo("Last login updated successfully", logrus.Fields{
		"username": user.Username,
		"id":       user.ID,
	})
	return nil
}

// CheckUserPassword checks if the provided password matches the user's password hash.
func CheckUserPassword(db *gorm.DB, password string, username string) (*models.User, bool, error) {
	utils.LogInfo("Validating user credentials", logrus.Fields{
		"username": username,
	})

	user, err := GetUserByUsername(db, username)
	if err != nil {
		utils.LogWarn("Invalid username", logrus.Fields{
			"username": username,
		})
		return nil, false, errors.New("invalid username")
	}

	if !security.CheckPasswordHash(password, user.PasswordHash) {
		utils.LogWarn("Invalid password", logrus.Fields{
			"username": username,
		})
		return user, false, errors.New("invalid password")
	}

	utils.LogInfo("User credentials validated successfully", logrus.Fields{
		"username": username,
	})
	return user, true, nil
}
