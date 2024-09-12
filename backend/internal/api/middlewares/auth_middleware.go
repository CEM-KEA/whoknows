package middlewares

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type contextKey string

const UserKey contextKey = "userID"

// AuthMiddleware is a middleware function for checking the presence of a valid JWT token in the request headers.
// It expects a function that can validate the JWT token and return the claims if the token is valid.
// The user ID is extracted from the claims and added to the request context.
// If the token is invalid or missing, it returns a 401 Unauthorized response.
func AuthMiddleware(validateJWT func(token string) (map[string]interface{}, error)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract the token from the Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header is required", http.StatusUnauthorized)
				return
			}

			// The token is expected to be in the format "Bearer <token>"
			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
				return
			}

			token := tokenParts[1]

			// Validate the token
			claims, err := validateJWT(token)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Add user information to the request context
			userID := claims["sub"].(string) // Assuming "sub" contains the user ID
			ctx := context.WithValue(r.Context(), UserKey, userID)

			// Pass the request with the context to the next handler
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserIDFromContext is a helper function to retrieve the user ID from the context.
func GetUserIDFromContext(ctx context.Context) (uint, error) {
	if userIDStr, ok := ctx.Value(UserKey).(string); ok {
		userIDUint, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			return 0, err
		}
		return uint(userIDUint), nil
	}
	return 0, errors.New("user ID not found in context")
}
