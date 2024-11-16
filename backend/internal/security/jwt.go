package security

import (
	"time"

	"github.com/CEM-KEA/whoknows/backend/internal/config"
	"github.com/CEM-KEA/whoknows/backend/internal/models"
	"github.com/CEM-KEA/whoknows/backend/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// GenerateJWT generates a JWT token for a given user ID and username
func GenerateJWT(userID uint, username string) (string, error) {
	utils.LogInfo("Generating JWT token", logrus.Fields{"userID": userID, "username": username})

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
		utils.LogError(err, "Failed to sign token", nil)
		return "", err
	}

	utils.LogInfo("Token signed successfully", logrus.Fields{"userID": userID, "username": username})
	return tokenString, nil
}

// GenerateJWTWithCustomExpiration generates a JWT token with a custom expiration time
func GenerateJWTWithCustomExpiration(userID uint, username string, expTime time.Time) (string, error) {
	utils.LogInfo("Generating JWT with custom expiration", logrus.Fields{
		"userID": userID, "username": username, "expTime": expTime,
	})

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
		utils.LogError(err, "Failed to sign token", nil)
		return "", err
	}

	utils.LogInfo("Token signed successfully", logrus.Fields{"userID": userID, "username": username})
	return tokenString, nil
}

// ValidateJWT validates the given JWT token and returns the claims if the token is valid
func ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	utils.LogInfo("Validating JWT token", logrus.Fields{})

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWT.Secret), nil
	})
	if err != nil {
		utils.LogError(err, "Failed to parse token", nil)
		return nil, errors.Wrap(err, "failed to parse token")
	}

	if !token.Valid {
		utils.LogInfo("Token is invalid", nil)
		return nil, errors.New("invalid token")
	}

	utils.LogInfo("Token validated successfully", nil)
	return claims, nil
}

// ValidateJWTRevoked checks if the token is revoked
func ValidateJWTRevoked(db *gorm.DB, jwt string) error {
	utils.LogInfo("Checking if JWT is revoked", logrus.Fields{"jwt": jwt})

	var jwtModel models.JWT
	err := db.Where("token = ?", jwt).First(&jwtModel).Error
	if err != nil {
		utils.LogError(err, "Failed to query token", logrus.Fields{"jwt": jwt})
		return errors.Wrap(err, "failed to query token")
	}

	if jwtModel.RevokedAt != nil {
		utils.LogInfo("Token is revoked", logrus.Fields{"jwt": jwt})
		return errors.New("token is revoked")
	}

	return nil
}

// RevokeJWT revokes a given JWT token
func RevokeJWT(db *gorm.DB, jwt string) error {
	utils.LogInfo("Revoking JWT token", logrus.Fields{"jwt": jwt})

	var jwtModel models.JWT
	err := db.Where("token = ?", jwt).First(&jwtModel).Error
	if err != nil {
		utils.LogError(err, "Failed to query token", logrus.Fields{"jwt": jwt})
		return errors.Wrap(err, "failed to query token")
	}

	jwtModel.RevokedAt = &time.Time{}
	err = db.Save(&jwtModel).Error
	if err != nil {
		utils.LogError(err, "Failed to update token", logrus.Fields{"jwt": jwt})
		return errors.Wrap(err, "failed to update token")
	}

	utils.LogInfo("Token revoked successfully", logrus.Fields{"jwt": jwt})
	return nil
}
