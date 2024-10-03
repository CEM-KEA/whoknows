package integration_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CEM-KEA/whoknows/backend/internal/api"
	"github.com/CEM-KEA/whoknows/backend/test/helpers"
	"github.com/stretchr/testify/assert"
)

func TestSearchIntegration(t *testing.T) {
	helpers.SetupTestDB(t)

	router := api.NewRouter()

	tests := []struct {
		name            string
		query           string
		expectedStatus  int
		expectedContent string
	}{
		{
			name:            "Valid Search Query",
			query:           "/api/search?q=programming&language=en",
			expectedStatus:  http.StatusOK,
			expectedContent: "Go Programming",
		},
		{
			name:            "Language Filtered Search",
			query:           "/api/search?q=Guide&language=da",
			expectedStatus:  http.StatusOK,
			expectedContent: "Danish Guide",
		},
		{
			name:            "Missing Query Parameter",
			query:           "/api/search",
			expectedStatus:  http.StatusBadRequest,
			expectedContent: "Search query (q) is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", tt.query, nil)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.expectedContent)
		})
	}
}
