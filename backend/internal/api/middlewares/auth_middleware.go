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

// AuthMiddleware validates the JWT token and adds the user ID to the context.
func AuthMiddleware(db *gorm.DB, validateJWT func(token string) (map[string]interface{}, error)) func(http.Handler) http.Handler {
	utils.LogInfo("Setting up auth middleware", nil)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract and validate the token
			token, err := extractToken(r)
			if err != nil {
				utils.LogWarn("Failed to extract token", logrus.Fields{"error": err.Error()})
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			// Validate the token and extract claims
			claims, err := validateJWT(token)
			if err != nil {
				utils.LogWarn("Invalid token", logrus.Fields{"error": err.Error()})
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Extract user ID from claims
			userID, err := extractUserID(claims)
			if err != nil {
				utils.LogWarn("Invalid user ID in token", logrus.Fields{"error": err.Error()})
				http.Error(w, "Invalid user ID", http.StatusUnauthorized)
				return
			}

			// Verify user exists in the database
			_, err = services.GetUserByID(db, uint(userID))
			if err != nil {
				utils.LogWarn("User not found", logrus.Fields{"userID": userID, "error": err.Error()})
				http.Error(w, "User not found", http.StatusUnauthorized)
				return
			}

			// Add user ID to the context and proceed
			ctx := context.WithValue(r.Context(), UserKey, uint(userID))
			utils.LogInfo("User authenticated successfully", logrus.Fields{"userID": userID})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// extractToken retrieves the JWT token from the Authorization header.
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

// extractUserID retrieves the user ID from JWT claims.
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

// GetUserIDFromContext retrieves the user ID from the request context.
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
