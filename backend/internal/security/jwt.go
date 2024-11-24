package security

import (
	"time"

	"github.com/CEM-KEA/whoknows/backend/internal/config"
	"github.com/CEM-KEA/whoknows/backend/internal/models"
	"github.com/CEM-KEA/whoknows/backend/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// GenerateJWT generates a JSON Web Token (JWT) for a given user ID and username.
// The token includes claims such as issuer, subject, audience, username, role, issued at, and expiration time.
// The token is signed using the HS256 signing method and a secret key from the application configuration.
//
// Parameters:
//   - userID: The unique identifier of the user.
//   - username: The username of the user.
//
// Returns:
//   - A signed JWT token as a string.
//   - An error if there is a failure in signing the token.
func GenerateJWT(userID uint, username string) (string, error) {
	utils.LogInfo("Starting JWT generation process", nil)

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
		utils.LogError(err, "Failed to sign JWT", nil)
		return "", err
	}

	utils.LogInfo("JWT generation completed successfully", nil)
	return tokenString, nil
}

// GenerateJWTWithCustomExpiration generates a JWT token with a custom expiration time.
// It takes the user ID, username, and expiration time as parameters and returns the signed JWT token string or an error.
//
// Parameters:
//   - userID: The unique identifier of the user.
//   - username: The username of the user.
//   - expTime: The custom expiration time for the JWT token.
//
// Returns:
//   - string: The signed JWT token string.
//   - error: An error if the token signing process fails.
func GenerateJWTWithCustomExpiration(userID uint, username string, expTime time.Time) (string, error) {
	utils.LogInfo("Starting JWT generation with custom expiration", nil)

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
		utils.LogError(err, "Failed to sign JWT with custom expiration", nil)
		return "", err
	}

	utils.LogInfo("JWT with custom expiration generated successfully", nil)
	return tokenString, nil
}

// ValidateJWT validates a given JWT token string and returns the claims if the token is valid.
// It logs the validation process and any errors encountered.
//
// Parameters:
//   - tokenString: The JWT token string to be validated.
//
// Returns:
//   - jwt.MapClaims: The claims extracted from the token if it is valid.
//   - error: An error if the token is invalid or if there was an issue during parsing.
func ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	utils.LogInfo("Starting JWT validation", nil)

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWT.Secret), nil
	})
	if err != nil {
		utils.LogError(err, "JWT validation failed during parsing", nil)
		return nil, errors.Wrap(err, "failed to parse JWT")
	}

	if !token.Valid {
		utils.LogInfo("JWT validation failed: token is invalid", nil)
		return nil, errors.New("invalid JWT token")
	}

	utils.LogInfo("JWT validation succeeded", nil)
	return claims, nil
}

// ValidateJWTRevoked checks if the token is revoked
// ValidateJWTRevoked checks if a given JWT token is revoked by querying the database.
// It logs the process of checking and any errors encountered during the query.
// If the token is found and is revoked, it returns an error indicating the token is revoked.
//
// Parameters:
//   - db: A pointer to the gorm.DB instance used to query the database.
//   - jwt: The JWT token string to be validated.
//
// Returns:
//   - error: An error if the token is revoked or if there is a failure in querying the database.
func ValidateJWTRevoked(db *gorm.DB, jwt string) error {
	utils.LogInfo("Checking if JWT is revoked", nil)

	var jwtModel models.JWT
	err := db.Where("token = ?", jwt).First(&jwtModel).Error
	if err != nil {
		utils.LogError(err, "Database query for JWT revocation check failed", nil)
		return errors.Wrap(err, "failed to query token revocation status")
	}

	if jwtModel.RevokedAt != nil {
		utils.LogInfo("JWT is revoked", nil)
		return errors.New("JWT token is revoked")
	}

	utils.LogInfo("JWT revocation check completed: token is not revoked", nil)
	return nil
}

// RevokeJWT revokes a given JWT token by marking it as revoked in the database.
//
// Parameters:
//   - db: A gorm.DB instance representing the database connection.
//   - jwt: A string representing the JWT token to be revoked.
//
// Returns:
//   - error: An error object if any error occurs during the process, otherwise nil.
//
// The function performs the following steps:
//  1. Logs the intention to revoke the JWT token.
//  2. Queries the database to find the JWT token.
//  3. If the token is found, it updates the token's revoked status.
//  4. Logs the success or failure of the revocation process.
func RevokeJWT(db *gorm.DB, jwt string) error {
	utils.LogInfo("Starting JWT revocation process", nil)

	var jwtModel models.JWT
	err := db.Where("token = ?", jwt).First(&jwtModel).Error
	if err != nil {
		utils.LogError(err, "Database query for JWT revocation failed", nil)
		return errors.Wrap(err, "failed to query token for revocation")
	}

	jwtModel.RevokedAt = &time.Time{}
	err = db.Save(&jwtModel).Error
	if err != nil {
		utils.LogError(err, "Failed to update JWT revocation status", nil)
		return errors.Wrap(err, "failed to update token revocation status")
	}

	utils.LogInfo("JWT revoked successfully", nil)
	return nil
}
