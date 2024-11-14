package middlewares

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/CEM-KEA/whoknows/backend/internal/services"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type contextKey string

const UserKey contextKey = "userID"

// AuthMiddleware checks for a valid JWT token and adds the user ID to the context.
func AuthMiddleware(db *gorm.DB, validateJWT func(token string) (map[string]interface{}, error)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract and validate the token
			token, err := extractToken(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			// Validate the token and get claims
			claims, err := validateJWT(token)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Extract user ID from claims
			userID, err := extractUserID(claims)
			if err != nil {
				http.Error(w, "Invalid user ID", http.StatusUnauthorized)
				return
			}

			// Get user from the database
			_, err = services.GetUserByID(db, uint(userID))
			if err != nil {
				http.Error(w, "User not found", http.StatusUnauthorized)
				return
			}

			// Add user ID to the request context
			ctx := context.WithValue(r.Context(), UserKey, uint(userID))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// extractToken retrieves the JWT from the Authorization header.
func extractToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("Authorization header is required")
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return "", errors.New("Invalid Authorization header format")
	}

	return tokenParts[1], nil
}

// extractUserID retrieves the user ID from JWT claims.
func extractUserID(claims map[string]interface{}) (uint64, error) {
	userIDStr, ok := claims["sub"].(string)
	if !ok {
		return 0, errors.New("user ID not found in token")
	}
	return strconv.ParseUint(userIDStr, 10, 32)
}

// GetUserIDFromContext is a helper function to retrieve the user ID from the context.
func GetUserIDFromContext(ctx context.Context) (uint, error) {
	userID, ok := ctx.Value(UserKey).(uint)
	if !ok {
		return 0, errors.New("user ID not found in context")
	}
	return userID, nil
}
