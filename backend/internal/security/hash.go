package security

import (
	"github.com/CEM-KEA/whoknows/backend/internal/utils"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"github.com/sirupsen/logrus"
)

// HashPassword hashes the password
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

// CheckPasswordHash checks the password hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		utils.LogWarn("Password hash comparison failed", logrus.Fields{"error": err.Error()})
		return false
	}

	utils.LogInfo("Password hash comparison successful", logrus.Fields{})
	return true
}
