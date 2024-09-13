package security

import (
	"fmt"
	"time"

	"github.com/CEM-KEA/whoknows/backend/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

var jwtKey = []byte(config.AppConfig.JWT.Secret)

// GenerateJWT generates a JWT token for a given user ID and email
func GenerateJWT(userID uint, email string) (string, error) {
	claims := jwt.MapClaims{
		"iss":   "whoknows",
		"sub":   userID,
		"aud":   "whoknows",
		"email": email,
		"role":  "user",
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := config.AppConfig.JWT.Secret

	// Debug: Print the secret key used for signing
	fmt.Printf("Secret Key in GenerateJWT: %s\n", secret)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GenerateJWTWithCustomExpiration generates a JWT token for a given user ID and email, with a custom expiration time
func GenerateJWTWithCustomExpiration(userID uint, email string, expTime time.Time) (string, error) {
	claims := jwt.MapClaims{
		"iss":   "whoknows",
		"sub":   userID,
		"aud":   "whoknows",
		"email": email,
		"role":  "user",
		"iat":   time.Now().Unix(),
		"exp":   expTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := config.AppConfig.JWT.Secret

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateJWT validates the given JWT token and returns the claims if the token is valid.
func ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	// Create a new map to store the claims
	claims := jwt.MapClaims{}

	// Parse the token with the claims
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	// Check if there was an error parsing the token
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse token")
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	fmt.Println("JWT token validated successfully")

	return claims, nil
}

// TODO: Implement a function to refresh JWT tokens
// TODO: Implement a function to revoke JWT tokens
// TODO: Implement a function to blacklist JWT tokens
// TODO: Implement a function to check if a JWT token is blacklisted
