package middlewares

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/CEM-KEA/whoknows/backend/internal/services"
	"github.com/CEM-KEA/whoknows/backend/internal/utils"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type contextKey string

const UserKey contextKey = "userID"

// AuthMiddleware is a middleware function that handles authentication for incoming HTTP requests.
// It extracts and validates a JWT token from the request, verifies the user exists in the database,
// and adds the user ID to the request context if authentication is successful.
//
// Parameters:
//   - db: A *gorm.DB instance for database operations.
//   - validateJWT: A function that takes a JWT token string and returns the token claims as a map and an error.
//
// Returns:
//   A middleware function that wraps an http.Handler and performs authentication.
func AuthMiddleware(db *gorm.DB, validateJWT func(token string) (map[string]interface{}, error)) func(http.Handler) http.Handler {
	utils.LogInfo("Setting up auth middleware", nil)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := extractToken(r)
			if err != nil {
				utils.LogWarn("Failed to extract token", logrus.Fields{"error": err.Error()})
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			claims, err := validateJWT(token)
			if err != nil {
				utils.LogWarn("Invalid token", logrus.Fields{"error": err.Error()})
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			userID, err := extractUserID(claims)
			if err != nil {
				utils.LogWarn("Invalid user ID in token", logrus.Fields{"error": err.Error()})
				http.Error(w, "Invalid user ID", http.StatusUnauthorized)
				return
			}

			_, err = services.GetUserByID(db, uint(userID))
			if err != nil {
				utils.LogWarn("User not found", logrus.Fields{"userID": userID, "error": err.Error()})
				http.Error(w, "User not found", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserKey, uint(userID))
			utils.LogInfo("User authenticated successfully", logrus.Fields{"userID": userID})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// extractToken extracts the Bearer token from the Authorization header of an HTTP request.
// It logs the process of extracting the token and returns an error if the Authorization header
// is missing or if it is not in the correct format.
//
// Parameters:
//   - r: The HTTP request from which to extract the token.
//
// Returns:
//   - string: The extracted token if successful.
//   - error: An error if the Authorization header is missing or incorrectly formatted.
func extractToken(r *http.Request) (string, error) {
	utils.LogInfo("Extracting token from request", nil)

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		utils.LogWarn("Authorization header is missing", nil)
		return "", errors.New("Authorization header is required")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		utils.LogWarn("Invalid Authorization header format", logrus.Fields{"authHeader": authHeader})
		return "", errors.New("Invalid Authorization header format")
	}

	utils.LogInfo("Token extracted successfully", nil)
	return parts[1], nil
}


// extractUserID extracts the user ID from the given JWT claims map.
// It expects the user ID to be stored under the "sub" key as a string.
// If the user ID is found and successfully parsed as a uint64, it returns the user ID.
// Otherwise, it returns an error indicating the failure reason.
//
// Parameters:
//   - claims: A map containing JWT claims.
//
// Returns:
//   - uint64: The extracted user ID.
//   - error: An error if the user ID is not found or cannot be parsed.
func extractUserID(claims map[string]interface{}) (uint64, error) {
	utils.LogInfo("Extracting user ID from claims", nil)

	userIDStr, ok := claims["sub"].(string)
	if !ok {
		utils.LogWarn("User ID not found in claims", nil)
		return 0, errors.New("User ID not found in token claims")
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		utils.LogWarn("Failed to parse user ID", logrus.Fields{"userIDStr": userIDStr, "error": err.Error()})
		return 0, errors.Wrap(err, "Failed to parse user ID")
	}

	utils.LogInfo("User ID extracted successfully", logrus.Fields{"userID": userID})
	return userID, nil
}


// GetUserIDFromContext retrieves the user ID from the given context.
// It expects the user ID to be stored in the context with the key UserKey.
// If the user ID is found, it returns the user ID and a nil error.
// If the user ID is not found, it returns 0 and an error indicating that the user ID was not found.
//
// Parameters:
//   - ctx: The context from which to retrieve the user ID.
//
// Returns:
//   - uint: The user ID retrieved from the context.
//   - error: An error if the user ID is not found in the context.
func GetUserIDFromContext(ctx context.Context) (uint, error) {
	utils.LogInfo("Getting user ID from context", nil)

	userID, ok := ctx.Value(UserKey).(uint)
	if !ok {
		utils.LogWarn("User ID not found in context", nil)
		return 0, errors.New("User ID not found in context")
	}

	utils.LogInfo("User ID retrieved from context", logrus.Fields{"userID": userID})
	return userID, nil
}
