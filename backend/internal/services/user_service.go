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

// It logs the process of creating the user, including the username and email.
// If the creation fails, it logs the error and returns it.
//
// Parameters:
//   - db: A gorm.DB instance representing the database connection.
//   - user: A pointer to a models.User instance representing the user to be created.
//
// Returns:
//   - error: An error if the user creation fails, otherwise nil.
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

// UpdateUser updates the given user in the database.
// It logs the process of updating the user and returns an error if the update fails.
//
// Parameters:
//   - db: A pointer to the gorm.DB instance.
//   - user: A pointer to the models.User instance to be updated.
//
// Returns:
//   - error: An error if the update operation fails, otherwise nil.
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

// GetUserByUsername retrieves a user from the database by their username.
// It logs the process of retrieving the user and any errors that occur.
//
// Parameters:
//   - db: A pointer to the gorm.DB instance used to interact with the database.
//   - username: The username of the user to retrieve.
//
// Returns:
//   - A pointer to the retrieved models.User instance, or nil if an error occurred.
//   - An error if the user could not be retrieved, or nil if the operation was successful.
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


// GetUserByID retrieves a user from the database by their ID.
// It logs the process of retrieving the user and any errors encountered.
// 
// Parameters:
//   - db: A pointer to the gorm.DB instance used to interact with the database.
//   - userID: The ID of the user to retrieve.
//
// Returns:
//   - A pointer to the retrieved models.User instance, or nil if an error occurred.
//   - An error if the user could not be found or another issue occurred during retrieval.
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


// UpdateLastLogin updates the LastLogin field of the given user to the current time
// and saves the changes to the database.
//
// Parameters:
//   - db: A pointer to the gorm.DB instance used to interact with the database.
//   - user: A pointer to the models.User instance whose LastLogin field needs to be updated.
//
// Returns:
//   - error: An error object if the update operation fails, otherwise nil.
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


// CheckUserPassword validates the provided username and password against the database.
// It returns the user object, a boolean indicating whether the credentials are valid, and an error if any occurred.
//
// Parameters:
//   - db: A pointer to the gorm.DB instance for database operations.
//   - password: The password string to be validated.
//   - username: The username string to be validated.
//
// Returns:
//   - *models.User: A pointer to the User object if the credentials are valid.
//   - bool: A boolean indicating whether the credentials are valid.
//   - error: An error object if any error occurred during the validation process.
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
