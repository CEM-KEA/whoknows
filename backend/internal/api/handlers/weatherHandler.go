package handlers

import (
	"encoding/json"
	"net/http"
)

// Structure of the weather response
type WeatherResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Data       string `json:"data"` //Field for the weather data
}

// Handler for the weather API
func Weather(w http.ResponseWriter, r *http.Request) {

	//Retrieves weather data
	weatherData, err := getWeatherData()

	//Chackes if there's an error
	if err != nil {
		//Error response
		response := WeatherResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Retrieving weather data failed",
			Data:       "",
		}
		writeJsonResponse(w, response)
		return
	}
	//Successful response
	response := WeatherResponse{
		StatusCode: http.StatusOK,
		Message:    "Successfully retrieved weather data",
		Data:       weatherData,
	}
	writeJsonResponse(w, response)
}

// Writes and encodes the response for the client
func writeJsonResponse(w http.ResponseWriter, response WeatherResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)
	// Converts the response object to JSON, writes it to http.ResponseWrite,
	// and sends it to the client
	json.NewEncoder(w).Encode(response)
}

// Get data from database
func getWeatherData() (string, error) {
	return "Sample Weather Data", nil //Test data
}
