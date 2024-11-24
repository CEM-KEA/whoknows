package services

import (
	"time"

	"github.com/CEM-KEA/whoknows/backend/internal/models"
	"github.com/CEM-KEA/whoknows/backend/internal/security"
	"github.com/CEM-KEA/whoknows/backend/internal/utils"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)


// CreateUser creates a new user in the database.
// It logs the creation attempt and any errors that occur during the process.
//
// Parameters:
//   - db: A pointer to the gorm.DB instance used to interact with the database.
//   - user: A pointer to the models.User instance representing the user to be created.
//
// Returns:
//   - error: An error if the user creation fails, otherwise nil.
func CreateUser(db *gorm.DB, user *models.User) error {
	utils.LogInfo("Creating new user", utils.SanitizeFields(map[string]interface{}{
		"username": user.Username,
		"email":    user.Email,
	}))

	err := db.Create(user).Error
	if err != nil {
		utils.LogError(err, "Failed to create user", utils.SanitizeFields(map[string]interface{}{
			"username": user.Username,
			"email":    user.Email,
		}))
		return errors.Wrap(err, "failed to create user")
	}

	return nil
}


// UpdateUser updates the given user in the database.
// It logs the update operation and returns an error if the update fails.
//
// Parameters:
//   - db: A pointer to the gorm.DB instance.
//   - user: A pointer to the models.User instance to be updated.
//
// Returns:
//   - error: An error if the update operation fails, otherwise nil.
func UpdateUser(db *gorm.DB, user *models.User) error {
	utils.LogInfo("Updating user", map[string]interface{}{
		"id": user.ID,
	})

	err := db.Save(user).Error
	if err != nil {
		utils.LogError(err, "Failed to update user", map[string]interface{}{
			"id": user.ID,
		})
		return errors.Wrap(err, "failed to update user")
	}

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
//   - A pointer to the retrieved models.User instance if found.
//   - An error if the user is not found or if there is an issue with the database query.
func GetUserByUsername(db *gorm.DB, username string) (*models.User, error) {
	utils.LogInfo("Retrieving user by username", map[string]interface{}{
		"username": username,
	})

	user := &models.User{}
	err := db.Where("username = ?", username).First(user).Error
	if err != nil {
		utils.LogError(err, "Failed to retrieve user by username", map[string]interface{}{
			"username": username,
		})
		return nil, errors.New("user not found")
	}

	return user, nil
}


// GetUserByID retrieves a user from the database by their ID.
// It logs the process of retrieving the user and any errors that occur.
//
// Parameters:
//   - db: A pointer to the gorm.DB instance used to interact with the database.
//   - userID: The ID of the user to retrieve.
//
// Returns:
//   - A pointer to the retrieved models.User instance, or nil if an error occurred.
//   - An error if the user could not be found or if there was an issue with the database query.
func GetUserByID(db *gorm.DB, userID uint) (*models.User, error) {
	utils.LogInfo("Retrieving user by ID", map[string]interface{}{
		"userID": userID,
	})

	user := &models.User{}
	err := db.First(user, userID).Error
	if err != nil {
		utils.LogError(err, "Failed to retrieve user by ID", map[string]interface{}{
			"userID": userID,
		})
		return nil, errors.New("user not found")
	}

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
	utils.LogInfo("Updating last login for user", map[string]interface{}{
		"id": user.ID,
	})

	user.LastLogin = time.Now()
	err := db.Save(user).Error
	if err != nil {
		utils.LogError(err, "Failed to update last login", map[string]interface{}{
			"id": user.ID,
		})
		return errors.Wrap(err, "failed to update last login")
	}

	return nil
}


// CheckUserPassword validates the provided username and password against the stored credentials in the database.
// It logs the validation process and returns the user object if the credentials are valid.
//
// Parameters:
//   - db: A pointer to the gorm.DB instance for database operations.
//   - password: The password string provided by the user.
//   - username: The username string provided by the user.
//
// Returns:
//   - *models.User: A pointer to the User model if the credentials are valid.
//   - bool: A boolean indicating whether the credentials are valid.
//   - error: An error object if the credentials are invalid or if there is an issue during the process.
func CheckUserPassword(db *gorm.DB, password, username string) (*models.User, bool, error) {
	utils.LogInfo("Validating user credentials", map[string]interface{}{
		"username": username,
	})

	user, err := GetUserByUsername(db, username)
	if err != nil {
		utils.LogWarn("Invalid username", map[string]interface{}{
			"username": username,
		})
		return nil, false, errors.New("invalid username")
	}

	if !security.CheckPasswordHash(password, user.PasswordHash) {
		utils.LogWarn("Invalid password", map[string]interface{}{
			"username": username,
		})
		return user, false, errors.New("invalid password")
	}

	return user, true, nil
}
