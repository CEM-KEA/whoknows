package unit_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/CEM-KEA/whoknows/backend/internal/config"
	"github.com/CEM-KEA/whoknows/backend/internal/security"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateJWT(t *testing.T) {
	userID := uint(1)
	email := "test@example.com"

	tokenString, err := security.GenerateJWT(userID, email)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Debug: Print the secret key
		t.Logf("Secret Key: %s", config.AppConfig.JWT.Secret)
		return []byte(config.AppConfig.JWT.Secret), nil
	})
	assert.NoError(t, err)

	// Debug: Print the claims
	t.Logf("Claims: %+v", claims)

	assert.Equal(t, "whoknows", claims["iss"])
	assert.Equal(t, "whoknows", claims["aud"])
	assert.Equal(t, email, claims["email"])
	assert.Equal(t, "user", claims["role"])

	// Convert sub claim to string before assertion
	subClaim := fmt.Sprintf("%v", claims["sub"])
	assert.Equal(t, "1", subClaim)
}

func TestValidateJWT(t *testing.T) {
	// Generate a JWT
	userID := uint(1)
	token, err := security.GenerateJWT(userID, "test@example.com")
	assert.NoError(t, err)

	// Validate the JWT
	claims, err := security.ValidateJWT(token)
	assert.NoError(t, err)
	assert.NotNil(t, claims)

	// Convert uint to string for comparison
	assert.Equal(t, "1", fmt.Sprintf("%d", userID))
}

func TestValidateJWT_InvalidToken(t *testing.T) {
	invalidTokenString := "invalid.token.string"

	claims, err := security.ValidateJWT(invalidTokenString)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
	userID := uint(1)
	// Generate an expired JWT by setting exp to a past time
	expiredTime := time.Now().Add(-time.Hour) // Set to 1 hour in the past
	token, err := security.GenerateJWTWithCustomExpiration(userID, "test@example.com", expiredTime)
	assert.NoError(t, err)

	// Validate the expired JWT
	claims, err := security.ValidateJWT(token)

	// We expect an error due to the token being expired
	assert.Error(t, err, "An error is expected but got nil.")

	// Ensure no claims are returned
	assert.Nil(t, claims)
}
