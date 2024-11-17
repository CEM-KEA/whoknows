package unit_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/CEM-KEA/whoknows/backend/internal/api/middlewares"
	"github.com/CEM-KEA/whoknows/backend/internal/models"
	"github.com/CEM-KEA/whoknows/backend/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func init() {
	// Initialize logger for tests with debug level and text format
	utils.InitGlobalLogger("debug", "text")
}

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Auto migrate the user model
	err = db.AutoMigrate(&models.User{})
	require.NoError(t, err)

	return db
}

// createTestUser creates a test user in the database
func createTestUser(t *testing.T, db *gorm.DB) *models.User {
	user := &models.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
	}
	
	result := db.Create(user)
	require.NoError(t, result.Error)
	return user
}

func TestAuthMiddleware(t *testing.T) {
	// Setup test database
	db := setupTestDB(t)
	testUser := createTestUser(t, db)

	// Create mock JWT validator
	mockValidateJWT := func(token string) (map[string]interface{}, error) {
		switch token {
		case "valid_token":
			return map[string]interface{}{
				"sub": strconv.FormatUint(uint64(testUser.ID), 10),
			}, nil
		case "invalid_user_token":
			return map[string]interface{}{
				"sub": "999", // Non-existent user ID
			}, nil
		case "malformed_sub_token":
			return map[string]interface{}{
				"sub": "not_a_number",
			}, nil
		default:
			return nil, errors.New("invalid token")
		}
	}

	tests := []struct {
		name           string
		setupRequest   func() *http.Request
		expectedStatus int
		expectedUserID uint
		expectError    bool
		errorMessage   string
	}{
		{
			name: "Valid token and existing user",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				req.Header.Set("Authorization", "Bearer valid_token")
				return req
			},
			expectedStatus: http.StatusOK,
			expectedUserID: testUser.ID,
			expectError:    false,
		},
		{
			name: "Missing Authorization header",
			setupRequest: func() *http.Request {
				return httptest.NewRequest(http.MethodGet, "/", nil)
			},
			expectedStatus: http.StatusUnauthorized,
			expectError:    true,
			errorMessage:   "Authorization header is required",
		},
		{
			name: "Invalid Authorization format",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				req.Header.Set("Authorization", "InvalidFormat")
				return req
			},
			expectedStatus: http.StatusUnauthorized,
			expectError:    true,
			errorMessage:   "Invalid Authorization header format",
		},
		{
			name: "Invalid token",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				req.Header.Set("Authorization", "Bearer invalid_token")
				return req
			},
			expectedStatus: http.StatusUnauthorized,
			expectError:    true,
			errorMessage:   "Invalid token",
		},
		{
			name: "Valid token but non-existent user",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				req.Header.Set("Authorization", "Bearer invalid_user_token")
				return req
			},
			expectedStatus: http.StatusUnauthorized,
			expectError:    true,
			errorMessage:   "User not found",
		},
		{
			name: "Malformed user ID in token",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				req.Header.Set("Authorization", "Bearer malformed_sub_token")
				return req
			},
			expectedStatus: http.StatusUnauthorized,
			expectError:    true,
			errorMessage:   "Invalid user ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a response recorder
			rr := httptest.NewRecorder()

			// Create a test handler that will be wrapped by the middleware
			var capturedUserID uint
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if userID, err := middlewares.GetUserIDFromContext(r.Context()); err == nil {
					capturedUserID = userID
				}
				w.WriteHeader(http.StatusOK)
			})

			// Create and apply the middleware
			middleware := middlewares.AuthMiddleware(db, mockValidateJWT)
			handler := middleware(testHandler)

			// Execute the request
			handler.ServeHTTP(rr, tt.setupRequest())

			// Assert the response
			assert.Equal(t, tt.expectedStatus, rr.Code, "Status code mismatch")

			if tt.expectError {
				assert.Contains(t, rr.Body.String(), tt.errorMessage, "Error message mismatch")
			} else {
				assert.Equal(t, tt.expectedUserID, capturedUserID, "User ID mismatch")
			}
		})
	}
}

func TestGetUserIDFromContext(t *testing.T) {
	tests := []struct {
		name        string
		ctx         context.Context
		expectedID  uint
		expectError bool
	}{
		{
			name:        "Valid user ID in context",
			ctx:         context.WithValue(context.Background(), middlewares.UserKey, uint(1)),
			expectedID:  1,
			expectError: false,
		},
		{
			name:        "No user ID in context",
			ctx:         context.Background(),
			expectedID:  0,
			expectError: true,
		},
		{
			name:        "Invalid type in context",
			ctx:         context.WithValue(context.Background(), middlewares.UserKey, "invalid"),
			expectedID:  0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, err := middlewares.GetUserIDFromContext(tt.ctx)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, uint(0), userID)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, userID)
			}
		})
	}
}