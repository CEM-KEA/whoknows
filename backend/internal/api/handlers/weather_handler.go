package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/CEM-KEA/whoknows/backend/internal/cache"
	"github.com/CEM-KEA/whoknows/backend/internal/config"
)

// We cache the weather data for 1 hour to reduce amount of calls to the API,
// as we only have 1000 free calls per day
var WeatherCache = cache.NewCache()

type WeatherResponse struct {
	Data map[string]interface{} `json:"data"`
}

// WeatherResponse represents the weather response payload
// @Description Get weather information
// @Produce json
// @Success 200 {object} WeatherResponse
// @Failure 500 {string} string "Failed to fetch weather data"
// @Router /api/weather [get]
// handler for GET request to /api/weather
func WeatherHandler(w http.ResponseWriter, r *http.Request) {
	data, err := GetWeatherData();
	if err != nil {
		http.Error(w, "Failed to fetch weather data", http.StatusInternalServerError)
		return
	}

	response := WeatherResponse{
		Data: data,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// fetches current weather data for Copenhagen from the OpenWeatherMap API
func GetWeatherData() (weatherData map[string]interface{}, err error) {
    // check if data is in cache
    if cachedData, found := WeatherCache.Get("weatherData"); found {
        return cachedData.(map[string]interface{}), nil
    }

    apiKey := config.AppConfig.WeatherAPI.OpenWeatherAPIKey
    url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=Copenhagen&appid=%s", apiKey)
    res, err := http.Get(url)
    if err != nil {
        fmt.Println("Error fetching weather data:", err)
        return weatherData, err
    }
    defer res.Body.Close()

    err = json.NewDecoder(res.Body).Decode(&weatherData)
    if err != nil {
        fmt.Println("Error decoding weather data:", err)
        return weatherData, err
    }

    // store data in cache with expiration of 1 hour
    WeatherCache.Set("weatherData", weatherData, 1 * time.Hour)

    return weatherData, nil
}