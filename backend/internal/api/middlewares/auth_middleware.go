package middlewares

import (
	"context"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/CEM-KEA/whoknows/backend/internal/security"
)

type contextKey string

const userKey contextKey = "userID"

// AuthMiddleware is a middleware function for checking the presence of a valid JWT token in the request headers.
func AuthMiddleware(next http.Handler) http.Handler {
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
		claims, err := security.ValidateJWT(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add user information to the request context
		userID := claims["sub"].(string) // Assuming "sub" contains the user ID
		ctx := context.WithValue(r.Context(), userKey, userID)

		// Pass the request with the context to the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserIDFromContext is a helper function to retrieve the user ID from the context.
func GetUserIDFromContext(ctx context.Context) (uint, error) {
	if userIDStr, ok := ctx.Value(userKey).(string); ok {
		userIDUint, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			return 0, err
		}
		return uint(userIDUint), nil
	}
	return 0, errors.New("user ID not found in context")
}
