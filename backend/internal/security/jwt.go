package security

import (
	"fmt"
	"time"

	"github.com/CEM-KEA/whoknows/backend/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

var jwtKey = []byte(config.AppConfig.JWT.Secret)

// GenerateJWT generates a JWT token for a user with the given userID.
func GenerateJWT(userID uint, email string) (string, error) {

	// Create the JWT claims, which include the user ID, email, and role
	claims := jwt.MapClaims{
		"iss":   "whoknows",                                                                              			// Issuer
		"sub":   fmt.Sprintf("%d", userID),                                                                        	// Subject (user ID as string)
		"aud":   "whoknows",                                                                              			// Audience 
		"exp":   jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(config.AppConfig.JWT.Expiration))), 	// Expiration time
		"iat":   jwt.NewNumericDate(time.Now()),                                                                   	// Issued at time
		"email": email,																								// Custom claim: user email
		"role":  "user",                                                                                           	// Custom claim: user role
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	return token.SignedString(jwtKey)
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