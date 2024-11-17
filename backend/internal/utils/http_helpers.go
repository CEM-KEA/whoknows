package utils

import (
	"encoding/json"
	"net/http"
)

// WriteJSONError sends a JSON response with an error message and HTTP status code.
// It sets the "Content-Type" header to "application/json" and writes the provided
// status code and error message to the response.
//
// Parameters:
//   - w: The http.ResponseWriter to write the response to.
//   - message: The error message to include in the JSON response.
//   - status: The HTTP status code to set for the response.
func WriteJSONError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}


// JSONSuccess sends a JSON response with the given data and status code.
// It sets the "Content-Type" header to "application/json" and writes the
// provided status code to the response header before encoding the data
// as JSON and writing it to the response body.
//
// Parameters:
//   w:    The http.ResponseWriter to write the response to.
//   data: A map containing the data to be encoded as JSON.
//   status: The HTTP status code to be set in the response header.
func JSONSuccess(w http.ResponseWriter, data map[string]interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}