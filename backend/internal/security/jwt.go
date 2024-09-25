package security

import (
	"fmt"
	"time"

	"github.com/CEM-KEA/whoknows/backend/internal/config"
	"github.com/CEM-KEA/whoknows/backend/internal/database"
	"github.com/CEM-KEA/whoknows/backend/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// GenerateJWT generates a JWT token for a given user ID and username
func GenerateJWT(userID uint, username string) (string, error) {
	claims := jwt.MapClaims{
		"iss":      "whoknows",
		"sub":      userID,
		"aud":      "whoknows",
		"username": username,
		"role":     "user",
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := config.AppConfig.JWT.Secret

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	jwtModel := models.JWT{
		UserID:    userID,
		Token:     tokenString,
		ExpiresAt: time.Now().Add(time.Hour * 24),
		CreatedAt: time.Now(),
		RevokedAt: nil,
	}

	err = database.DB.Create(&jwtModel).Error
	if err != nil {
		return "", errors.Wrap(err, "failed to add jwt token to database")
	}

	return tokenString, nil
}

// GenerateJWTWithCustomExpiration generates a JWT token for a given user ID and username, with a custom expiration time
func GenerateJWTWithCustomExpiration(userID uint, username string, expTime time.Time) (string, error) {
	claims := jwt.MapClaims{
		"iss":      "whoknows",
		"sub":      userID,
		"aud":      "whoknows",
		"username": username,
		"role":     "user",
		"iat":      time.Now().Unix(),
		"exp":      expTime.Unix(),
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
	claims := jwt.MapClaims{}

	// Parse the token with the claims
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		secret := config.AppConfig.JWT.Secret
		// Debug: Print the secret key used for validation
		fmt.Printf("Secret Key in ValidateJWT: %s\n", secret)

		return []byte(secret), nil
	})

	// Check if there was an error parsing the token
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse token")
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Check if the token is revoked
	var jwtModel models.JWT
	db := database.DB
	if err := db.Where("token = ?", tokenString).First(&jwtModel).Error; err != nil {
		return nil, errors.Wrap(err, "failed to query token")
	}

	if jwtModel.RevokedAt != nil {
		return nil, errors.New("token has been revoked")
	}

	fmt.Println("JWT token validated successfully")

	return claims, nil
}

// RevokeJWT revokes a given JWT token in the database
func RevokeJWT(db *gorm.DB, jwt string) error {
	var jwtModel models.JWT
	if err := db.Where("token = ?", jwt).First(&jwtModel).Error; err != nil {
		return errors.Wrap(err, "failed to query token")
	}

	jwtModel.RevokedAt = &time.Time{}
	if err := db.Save(&jwtModel).Error; err != nil {
		return errors.Wrap(err, "failed to update token")
	}

	return nil
}
