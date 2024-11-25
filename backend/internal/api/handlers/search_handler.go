package handlers

import (
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
//
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
	q := utils.SanitizeValue(r.URL.Query().Get("q"))
	language := utils.SanitizeValue(r.URL.Query().Get("language"))

	if q == "" {
		utils.LogWarn("Search query validation failed", nil)
		utils.WriteJSONError(w, "Search query (q) is required", http.StatusBadRequest)
		return
	}

	searchLog := models.SearchLog{Query: q}
	if err := services.CreateSearchLog(database.DB, &searchLog); err != nil {
		utils.LogError(err, "Failed to log search query", nil)
		utils.WriteJSONError(w, "Failed to log search query", http.StatusInternalServerError)
		return
	}

	queryType := "generic"
	if language != "" {
		queryType = "language_filtered"
	}
	utils.IncrementSearchQueries(queryType)

	var pages []models.Page
	query := database.DB.Where("content LIKE ?", "%"+q+"%").Order("title ASC")
	if language != "" {
		query = query.Where("language = ?", language)
	}
	if err := query.Find(&pages).Error; err != nil {
		utils.LogError(err, "Search query execution failed", nil)
		utils.WriteJSONError(w, "Search query failed", http.StatusInternalServerError)
		return
	}

	response := SearchResponse{
		Data: make([]map[string]interface{}, len(pages)),
	}
	for i, page := range pages {
		// get the content of the page around the first search match
		content := utils.GetContentAroundMatch(page.Content, q)
		response.Data[i] = map[string]interface{}{
			"id":       page.ID,
			"content":  content,
			"language": page.Language,
			"title":    page.Title,
			"url":      page.Url,
		}
	}

	utils.JSONSuccess(w, map[string]interface{}{
		"status":  "success",
		"data":    response.Data,
		"results": len(response.Data),
	}, http.StatusOK)

	utils.LogInfo("Search query completed successfully", logrus.Fields{
		"query":    q,
		"language": language,
		"results":  len(pages),
	})
}
