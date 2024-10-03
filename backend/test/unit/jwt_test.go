package unit_test

import (
	"testing"
	"time"

	"github.com/CEM-KEA/whoknows/backend/internal/security"
	"github.com/CEM-KEA/whoknows/backend/test/helpers"
	"github.com/stretchr/testify/assert"
)

const (
	testUsername = "testuser"
)

// TestGenerateAndValidateJWT tests the generation and validation of a JWT
func TestGenerateAndValidateJWT(t *testing.T) {
	// Setup test environment and database
	helpers.SetupTestDB(t)

	// Generate a valid JWT token
	token, err := security.GenerateJWT(1, testUsername)
	assert.NoError(t, err)

	// Validate the JWT token
	claims, err := security.ValidateJWT(token)
	assert.NoError(t, err)
	assert.NotNil(t, claims)

	// Ensure the expected claims are present
	assert.Equal(t, float64(1), claims["sub"]) // JWT stores numbers as float64
	assert.Equal(t, testUsername, claims["username"])
}

// TestValidateJWTInvalidToken tests JWT validation with an invalid token
func TestValidateJWTInvalidToken(t *testing.T) {
	// Setup test environment and database
	helpers.SetupTestDB(t)

	invalidTokenString := "invalid.token.string"

	// Validate an invalid JWT token
	claims, err := security.ValidateJWT(invalidTokenString)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

// TestValidateJWTExpiredToken tests JWT validation with an expired token
func TestValidateJWTExpiredToken(t *testing.T) {
	// Setup test environment and database
	helpers.SetupTestDB(t)

	// Generate an expired JWT by setting exp to a past time
	expiredTime := time.Now().Add(-time.Hour) // Set to 1 hour in the past
	token, err := security.GenerateJWTWithCustomExpiration(1, testUsername, expiredTime)
	assert.NoError(t, err)

	// Validate the expired JWT
	claims, err := security.ValidateJWT(token)

	// We expect an error due to the token being expired
	assert.Error(t, err, "An error is expected but got nil.")

	// Ensure no claims are returned
	assert.Nil(t, claims)
}
