package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CEM-KEA/whoknows/backend/internal/api"
	"github.com/CEM-KEA/whoknows/backend/test/helpers"
	"github.com/stretchr/testify/assert"
)

func TestRegisterIntegration(t *testing.T) {
	helpers.SetupTestDB(t)

	router := api.NewRouter()

	tests := []struct {
		name           string
		payload        map[string]string
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Valid Registration",
			payload: map[string]string{
				"username":  "newuser",
				"email":     "newuser@example.com",
				"password":  "password123",
				"password2": "password123",
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   "User created successfully",
		},
		{
			name: "Password Mismatch",
			payload: map[string]string{
				"username":  "newuser",
				"email":     "newuser@example.com",
				"password":  "password123",
				"password2": "password456",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Password confirmation does not match\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)

			req, _ := http.NewRequest("POST", "/api/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.expectedBody)
		})
	}
}
