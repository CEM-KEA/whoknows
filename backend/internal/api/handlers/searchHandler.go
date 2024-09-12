package handlers

import (
    "encoding/json"
    "net/http"
    "github.com/CEM-KEA/whoknows/backend/internal/database"
    "github.com/CEM-KEA/whoknows/backend/internal/models"
)

// SearchBody represents the structure of the request body
type SearchBody struct {
    Q        string  `json:"q"`
    Language *string `json:"language,omitempty"`
}

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
func Search(w http.ResponseWriter, r *http.Request) {
    var body SearchBody
    err := json.NewDecoder(r.Body).Decode(&body)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Validate the required search query (q)
    if body.Q == "" {
        msg := "Search query (q) is required"
        validationError := RequestValidationError{
            StatusCode: http.StatusBadRequest,
            Message:    &msg,
        }
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(validationError)
        return
    }

    // Perform the search using Gorm
    var pages []models.Page
    query := database.DB.Where("content LIKE ?", "%"+body.Q+"%")

    // Optional language filter
    if body.Language != nil {
        query = query.Where("language = ?", *body.Language)
    }

    err = query.Find(&pages).Error
    if err != nil {
        http.Error(w, "Failed to retrieve search results", http.StatusInternalServerError)
        return
    }

    // Prepare response
    response := SearchResponse{
        Data: make([]map[string]interface{}, len(pages)),
    }

    for i, page := range pages {
        response.Data[i] = map[string]interface{}{
            "id":        page.ID,
            "title":     page.Title,
            "url":       page.Url,
            "language":  page.Language,
            "content":   page.Content,
            "createdAt": page.CreatedAt,
            "updatedAt": page.UpdatedAt,
        }
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
