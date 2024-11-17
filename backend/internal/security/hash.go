package security

import (
	"github.com/CEM-KEA/whoknows/backend/internal/utils"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strconv"
)

var bcryptCost int


// init initializes the bcrypt cost value from the environment variable "BCRYPT_COST".
// If the environment variable is not set or contains an invalid value, it falls back to the default cost.
// The cost value must be within the range defined by bcrypt.MinCost and bcrypt.MaxCost.
func init() {
	cost, err := strconv.Atoi(os.Getenv("BCRYPT_COST"))
	if err != nil || cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = bcrypt.DefaultCost // Fallback to Go's recommended default cost
	}
	bcryptCost = cost
}

// HashPassword hashes the given password using bcrypt.
// It returns the hashed password as a string and an error if the hashing fails.
//
// Parameters:
//   - password: The plain text password to be hashed.
//
// Returns:
//   - A string representing the hashed password.
//   - An error if the password is empty or hashing fails.
func HashPassword(password string) (string, error) {
	if password == "" {
		utils.LogWarn("Password cannot be empty", nil)
		return "", errors.New("password cannot be empty")
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		utils.LogError(err, "Failed to hash password", nil)
		return "", errors.Wrap(err, "failed to hash password")
	}

	return string(hashedBytes), nil
}

// CheckPasswordHash compares a plaintext password with a hashed password.
// It returns true if they match, otherwise false.
//
// Parameters:
//   - password: The plaintext password to compare.
//   - hash: The hashed password to compare against.
//
// Returns:
//   - bool: true if the password matches the hash, false otherwise.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		utils.LogWarn("Password hash comparison failed", nil)
		return false
	}

	return true
}
