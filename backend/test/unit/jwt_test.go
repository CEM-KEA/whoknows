package unit_test

import (
	"os"
	"testing"
	"time"

	"github.com/CEM-KEA/whoknows/backend/internal/security"
	"github.com/stretchr/testify/assert"
)

// TestGenerateJWT generates and validates JWT
func TestGenerateJWT(t *testing.T) {
	os.Setenv("ENV_FILE_PATH", "../test.env")

	// Generate a valid JWT token
	token, err := security.GenerateJWT(1, "test@example.com")
	assert.NoError(t, err)

	// Validate the JWT token
	claims, err := security.ValidateJWT(token)
	assert.NoError(t, err)
	assert.NotNil(t, claims)

	// Ensure the expected claims are present
	assert.Equal(t, float64(1), claims["sub"]) // JWT stores numbers as float64
	assert.Equal(t, "test@example.com", claims["email"])

	os.Unsetenv("ENV_FILE_PATH")
}

func TestValidateJWT(t *testing.T) {
	os.Setenv("ENV_FILE_PATH", "../test.env")

	// Generate a valid JWT token
	token, err := security.GenerateJWT(1, "test@example.com")
	assert.NoError(t, err)

	// Validate the JWT token
	claims, err := security.ValidateJWT(token)
	assert.NoError(t, err)
	assert.NotNil(t, claims)

	// Ensure the expected claims are present
	assert.Equal(t, float64(1), claims["sub"]) // JWT stores numbers as float64
	assert.Equal(t, "test@example.com", claims["email"])

	os.Unsetenv("ENV_FILE_PATH")
}

func TestValidateJWT_InvalidToken(t *testing.T) {
	invalidTokenString := "invalid.token.string"

	claims, err := security.ValidateJWT(invalidTokenString)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
	os.Setenv("ENV_FILE_PATH", "../test.env")

	// Generate an expired JWT by setting exp to a past time
	expiredTime := time.Now().Add(-time.Hour) // Set to 1 hour in the past
	token, err := security.GenerateJWTWithCustomExpiration(1, "test@example.com", expiredTime)
	assert.NoError(t, err)

	// Validate the expired JWT
	claims, err := security.ValidateJWT(token)

	// We expect an error due to the token being expired
	assert.Error(t, err, "An error is expected but got nil.")

	// Ensure no claims are returned
	assert.Nil(t, claims)

	os.Unsetenv("ENV_FILE_PATH")
}
