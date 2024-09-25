package integration_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/CEM-KEA/whoknows/backend/internal/api"
	"github.com/CEM-KEA/whoknows/backend/internal/utils"
    "github.com/stretchr/testify/assert"
)

func TestLoginIntegration(t *testing.T) {
	utils.InitValidator() // Initialize the validator

    setupTestDB(t) // Set up the test database

    // Create a new router for the tests
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
            // Marshal the payload to JSON
            body, _ := json.Marshal(tt.payload)

            // Create a new POST request to /api/login
            req, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer(body))
            req.Header.Set("Content-Type", "application/json")

            // Create a ResponseRecorder
            rr := httptest.NewRecorder()

            // Serve the request
            router.ServeHTTP(rr, req)

            // Check the status code
            assert.Equal(t, tt.expectedStatus, rr.Code)

            // Check the response body contains the expected substring
            assert.Contains(t, rr.Body.String(), tt.expectedBody)
        })
    }
}
