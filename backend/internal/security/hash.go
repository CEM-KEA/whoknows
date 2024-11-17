package security

import (
	"github.com/CEM-KEA/whoknows/backend/internal/utils"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"github.com/sirupsen/logrus"
)

// HashPassword hashes the given password using bcrypt with a cost of 14.
// It logs the process of hashing the password and handles errors appropriately.
//
// Parameters:
//   - password: The plain text password to be hashed.
//
// Returns:
//   - A string representing the hashed password.
//   - An error if the password is empty or if hashing fails.
func HashPassword(password string) (string, error) {
	utils.LogInfo("Hashing password", logrus.Fields{})

	if password == "" {
		utils.LogWarn("Password cannot be empty", logrus.Fields{})
		return "", errors.New("password cannot be empty")
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		utils.LogError(err, "Failed to hash password", logrus.Fields{})
		return "", errors.Wrap(err, "failed to hash password")
	}

	utils.LogInfo("Password hashed successfully", logrus.Fields{})
	return string(bytes), nil
}

// CheckPasswordHash compares a plaintext password with a hashed password and returns
// true if they match, otherwise returns false. It logs a warning if the comparison fails
// and logs an info message if the comparison is successful.
//
// Parameters:
//   - password: the plaintext password to compare
//   - hash: the hashed password to compare against
//
// Returns:
//   - bool: true if the password matches the hash, false otherwise
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		utils.LogWarn("Password hash comparison failed", logrus.Fields{"error": err.Error()})
		return false
	}

	utils.LogInfo("Password hash comparison successful", logrus.Fields{})
	return true
}
