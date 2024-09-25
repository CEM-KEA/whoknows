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

func TestRegisterIntegration(t *testing.T) {
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
            // Marshal the payload to JSON
            body, _ := json.Marshal(tt.payload)

            // Create a new POST request to /api/register
            req, _ := http.NewRequest("POST", "/api/register", bytes.NewBuffer(body))
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
