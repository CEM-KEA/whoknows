package handlers

import (
	"encoding/json"
	"net/http"
)

// SearchResponse represents the structure of the search response
type WeatherResponse struct {
	Data map[string]interface{} `json:"data"`
}

// Search is the handler for the search API
func WeatherHandler(w http.ResponseWriter, r *http.Request) {
	
	data := make(map[string]interface{}, 0);
	data["test"] = "test"
	response := WeatherResponse{
		Data: data,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
