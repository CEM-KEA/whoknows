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
	helpers.SetupTestDB(t)

	token, err := security.GenerateJWT(1, testUsername)
	assert.NoError(t, err)

	claims, err := security.ValidateJWT(token)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, float64(1), claims["sub"])
	assert.Equal(t, testUsername, claims["username"])
}

// TestValidateJWTInvalidToken tests JWT validation with an invalid token
func TestValidateJWTInvalidToken(t *testing.T) {
	helpers.SetupTestDB(t)

	invalidTokenString := "invalid.token.string"

	claims, err := security.ValidateJWT(invalidTokenString)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

// TestValidateJWTExpiredToken tests JWT validation with an expired token
func TestValidateJWTExpiredToken(t *testing.T) {
	helpers.SetupTestDB(t)

	expiredTime := time.Now().Add(-time.Hour)
	token, err := security.GenerateJWTWithCustomExpiration(1, testUsername, expiredTime)
	assert.NoError(t, err)

	claims, err := security.ValidateJWT(token)

	assert.Error(t, err, "An error is expected but got nil.")
	assert.Nil(t, claims)
}
