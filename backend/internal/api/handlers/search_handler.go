package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/CEM-KEA/whoknows/backend/internal/database"
	"github.com/CEM-KEA/whoknows/backend/internal/models"
	"github.com/CEM-KEA/whoknows/backend/internal/services"
	"github.com/CEM-KEA/whoknows/backend/internal/utils"
	"github.com/sirupsen/logrus"
)

// SearchResponse represents the structure of the search response
type SearchResponse struct {
	Data []map[string]interface{} `json:"data"`
}

// RequestValidationError represents validation error details
type RequestValidationError struct {
	StatusCode int     `json:"statusCode"`
	Message    *string `json:"message,omitempty"`
}

// Search is the handler for the search API
//	@Description	Search for pages by content
//	@Produce		json
//	@Param			q			query		string	true	"Search query"
//	@Param			language	query		string	false	"Language filter"
//	@Success		200			{object}	SearchResponse
//	@Failure		400			{string}	string	"Search query (q) is required"
//	@Failure		500			{string}	string	"Search query failed"
//	@Router			/api/search [get]
func Search(w http.ResponseWriter, r *http.Request) {
	utils.LogInfo("Processing search request", nil)

	// Extract query parameters
	q := r.URL.Query().Get("q")
	language := r.URL.Query().Get("language")

	// Validate the required search query parameter
	if q == "" {
		msg := "Search query (q) is required"
		utils.LogWarn("Search query validation failed", logrus.Fields{
			"error": msg,
		})

		validationError := RequestValidationError{
			StatusCode: http.StatusBadRequest,
			Message:    &msg,
		}
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(validationError); err != nil {
			utils.LogError(err, "Failed to encode validation error", nil)
			http.Error(w, "Failed to encode validation error", http.StatusInternalServerError)
		}
		return
	}

	// Log the search query
	searchLog := models.SearchLog{Query: q}
	if err := services.CreateSearchLog(database.DB, &searchLog); err != nil {
		utils.LogError(err, "Failed to log search query", logrus.Fields{
			"query": q,
		})
		http.Error(w, "Failed to log search query", http.StatusInternalServerError)
		return
	}

	// Perform the search
	var pages []models.Page
	query := database.DB.Where("content LIKE ?", "%"+q+"%").Order("title ASC")
	if language != "" {
		query = query.Where("language = ?", language)
	}

	if err := query.Find(&pages).Error; err != nil {
		utils.LogError(err, "Search query execution failed", logrus.Fields{
			"query":    q,
			"language": language,
		})
		http.Error(w, "Search query failed", http.StatusInternalServerError)
		return
	}

	// Prepare the response
	response := SearchResponse{
		Data: make([]map[string]interface{}, len(pages)),
	}
	for i, page := range pages {
		response.Data[i] = map[string]interface{}{
			"id":       page.ID,
			"content":  page.Content,
			"language": page.Language,
			"title":    page.Title,
			"url":      page.Url,
		}
	}

	// Encode and send the response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		utils.LogError(err, "Failed to encode search response", logrus.Fields{
			"query":    q,
			"language": language,
		})
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	// Log success
	utils.LogInfo("Search query completed successfully", logrus.Fields{
		"query":    q,
		"language": language,
		"results":  len(pages),
	})
}