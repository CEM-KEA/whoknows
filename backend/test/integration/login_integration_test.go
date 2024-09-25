package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CEM-KEA/whoknows/backend/internal/api"
	"github.com/stretchr/testify/assert"
)

func TestLoginIntegration(t *testing.T) {
	setupTestDB(t)

	router := api.NewRouter()

	tests := []struct {
		name           string
		payload        map[string]string
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Valid Login",
			payload: map[string]string{
				"username": "testuser",
				"password": "password123",
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "token",
		},
		{
			name: "Invalid Password",
			payload: map[string]string{
				"username": "testuser",
				"password": "wrongpassword",
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid password",
		},
		{
			name: "Non-existent User",
			payload: map[string]string{
				"username": "nonexistent",
				"password": "password123",
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid username",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)

			req, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.expectedBody)
		})
	}
}
