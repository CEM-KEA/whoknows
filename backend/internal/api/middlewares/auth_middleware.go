package middlewares

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/CEM-KEA/whoknows/backend/internal/services"
	"github.com/CEM-KEA/whoknows/backend/internal/utils"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type contextKey string

const UserKey contextKey = "userID"

// AuthMiddleware is a middleware function for handling authentication.
// It extracts the JWT token from the request, validates it, and retrieves the user ID from the token claims.
// If the token is valid and the user exists in the database, the request is allowed to proceed with the user ID added to the context.
// Otherwise, it responds with an appropriate error message and status code.
//
// Parameters:
//   - db: A gorm.DB instance for database operations.
//   - validateJWT: A function that takes a JWT token string and returns the token claims as a map and an error, if any.
//
// Returns:
//   A middleware function that wraps an http.Handler and performs authentication.
func AuthMiddleware(db *gorm.DB, validateJWT func(token string) (map[string]interface{}, error)) func(http.Handler) http.Handler {
	utils.LogInfo("Setting up authentication middleware", nil)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := extractToken(r)
			if err != nil {
				utils.LogWarn("Failed to extract token", logSanitizedError(err))
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			claims, err := validateJWT(token)
			if err != nil {
				utils.LogWarn("Invalid token", logSanitizedError(err))
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			userID, err := extractUserID(claims)
			if err != nil {
				utils.LogWarn("Invalid user ID in token", logSanitizedError(err))
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			_, err = services.GetUserByID(db, uint(userID))
			if err != nil {
				utils.LogWarn("User not found", utils.SanitizeFields(map[string]interface{}{"userID": userID, "error": err.Error()}))
				http.Error(w, "User not found", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserKey, uint(userID))
			utils.LogInfo("User authenticated successfully", utils.SanitizeFields(map[string]interface{}{"userID": userID}))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// extractToken extracts the token from the Authorization header of the given HTTP request.
// It expects the Authorization header to be in the format "Bearer <token>".
// If the header is missing or not in the correct format, it returns an error.
//
// Parameters:
//   - r: The HTTP request from which to extract the token.
//
// Returns:
//   - string: The extracted token.
//   - error: An error if the Authorization header is missing or not in the correct format.
func extractToken(r *http.Request) (string, error) {
	utils.LogInfo("Extracting token from Authorization header", nil)

	authHeader := utils.SanitizeValue(r.Header.Get("Authorization"))
	if authHeader == "" {
		utils.LogWarn("Authorization header is missing", nil)
		return "", errors.New("Authorization header is required")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		utils.LogWarn("Invalid Authorization header format", nil)
		return "", errors.New("Invalid Authorization header format")
	}

	utils.LogInfo("Token extracted successfully", nil)
	return parts[1], nil
}

// extractUserID extracts the user ID from the given JWT claims.
// It expects the user ID to be stored under the "sub" key as a string.
// If the user ID is not found or cannot be parsed, it returns an error.
//
// Parameters:
//   - claims: A map containing JWT claims.
//
// Returns:
//   - uint64: The extracted user ID.
//   - error: An error if the user ID is not found or cannot be parsed.
func extractUserID(claims map[string]interface{}) (uint64, error) {
	utils.LogInfo("Extracting user ID from JWT claims", nil)

	userIDStr, ok := claims["sub"].(string)
	if !ok {
		utils.LogWarn("User ID not found in claims", nil)
		return 0, errors.New("User ID not found in token claims")
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		utils.LogWarn("Failed to parse user ID", utils.SanitizeFields(map[string]interface{}{"userIDStr": userIDStr, "error": err.Error()}))
		return 0, errors.Wrap(err, "Failed to parse user ID")
	}

	utils.LogInfo("User ID extracted successfully", utils.SanitizeFields(map[string]interface{}{"userID": userID}))
	return userID, nil
}

// GetUserIDFromContext retrieves the user ID from the given context.
// It expects the user ID to be stored in the context with the key UserKey.
// If the user ID is found, it returns the user ID and a nil error.
// If the user ID is not found, it logs a warning and returns an error.
//
// Parameters:
//
//	ctx (context.Context): The context from which to retrieve the user ID.
//
// Returns:
//
//	uint: The user ID retrieved from the context.
//	error: An error if the user ID is not found in the context.
func GetUserIDFromContext(ctx context.Context) (uint, error) {
	utils.LogInfo("Getting user ID from request context", nil)

	userID, ok := ctx.Value(UserKey).(uint)
	if !ok {
		utils.LogWarn("User ID not found in context", nil)
		return 0, errors.New("User ID not found in context")
	}

	utils.LogInfo("User ID retrieved successfully", utils.SanitizeFields(map[string]interface{}{"userID": userID}))
	return userID, nil
}

// logSanitizedError takes an error as input, sanitizes its message using utils.SanitizeValue,
// and returns a map containing the sanitized error message with the key "error".
func logSanitizedError(err error) map[string]interface{} {
	return map[string]interface{}{"error": utils.SanitizeValue(err.Error())}
}
