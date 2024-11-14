package unit_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/CEM-KEA/whoknows/backend/internal/api/middlewares"
	"github.com/CEM-KEA/whoknows/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Mock the ValidateJWT function
type MockJWTValidator struct {
	mock.Mock
}

func (m *MockJWTValidator) ValidateJWT(token string) (map[string]interface{}, error) {
	args := m.Called(token)

	if claims := args.Get(0); claims != nil {
		return claims.(map[string]interface{}), args.Error(1)
	}

	return nil, args.Error(1)
}

// Mock database setup
func setupTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto-migrate the User model
	if err := db.AutoMigrate(&models.User{}); err != nil {
		return nil, err
	}

	// Insert a test user
	db.Create(&models.User{ID: 123, Username: "testuser", Email: "testuser@example.com"})

	return db, nil
}

func TestAuthMiddleware(t *testing.T) {
	mockValidator := &MockJWTValidator{}

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		expectedUserID string
	}{
		{
			name:           "No Authorization header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid Authorization header format",
			authHeader:     "InvalidToken",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid token",
			authHeader:     "Bearer invalidtoken",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Valid token",
			authHeader:     "Bearer validtoken",
			expectedStatus: http.StatusOK,
			expectedUserID: "123",
		},
	}

	// Setup test database
	db, err := setupTestDB()
	assert.NoError(t, err)

	// Mock the ValidateJWT function
	mockValidator.On("ValidateJWT", "validtoken").Return(map[string]interface{}{"sub": "123"}, nil)
	mockValidator.On("ValidateJWT", "invalidtoken").Return(nil, errors.New("invalid token"))

	authMiddleware := middlewares.AuthMiddleware(db, mockValidator.ValidateJWT)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a request with the specified Authorization header
			req, err := http.NewRequest("GET", "/", nil)
			assert.NoError(t, err)
			req.Header.Set("Authorization", tt.authHeader)

			// Create a ResponseRecorder to capture the response
			rr := httptest.NewRecorder()

			// Create a dummy handler to pass to the middleware
			dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.expectedUserID != "" {
					userID, err := middlewares.GetUserIDFromContext(r.Context())
					assert.NoError(t, err)
					assert.Equal(t, tt.expectedUserID, strconv.FormatUint(uint64(userID), 10))
				}
				w.WriteHeader(http.StatusOK)
			})

			// Apply the middleware to the dummy handler
			handler := authMiddleware(dummyHandler)

			// Serve the request
			handler.ServeHTTP(rr, req)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
