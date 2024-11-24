package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/CEM-KEA/whoknows/backend/internal/cache"
	"github.com/CEM-KEA/whoknows/backend/internal/config"
	"github.com/CEM-KEA/whoknows/backend/internal/utils"
	"github.com/sirupsen/logrus"
)

// We cache the weather data for 1 hour to reduce amount of calls to the API,
// as we only have 1000 free calls per day
var WeatherCache = cache.NewCache()

type WeatherResponse struct {
	Data map[string]interface{} `json:"data"`
}

const (
	fetchDataError       = "Failed to fetch weather data"
	encodeResponseError  = "Failed to encode weather response"
	weatherDataCacheKey  = "weatherData"
	weatherDataCacheTime = 1 * time.Hour
)

// WeatherResponse represents the weather response payload
//	@Description	Get weather information
//	@Produce		json
//	@Success		200	{object}	WeatherResponse
//	@Failure		500	{string}	string	"Failed to fetch weather data"
//	@Router			/api/weather [get]
// handler for GET request to /api/weather
func WeatherHandler(w http.ResponseWriter, r *http.Request) {
	utils.LogInfo("Processing weather request", nil)
	data, err := GetWeatherData()
	if err != nil {
		utils.LogError(err, fetchDataError, nil)
		utils.WriteJSONError(w, fetchDataError, http.StatusInternalServerError)
		return
	}

	utils.JSONSuccess(w, map[string]interface{}{
		"status": "success",
		"data":   data,
	}, http.StatusOK)

	utils.LogInfo("Weather data fetched and response sent successfully", nil)
}


// GetWeatherData fetches weather data from the OpenWeather API for Copenhagen,
// logs the request and response process, and stores the fetched data in a cache.
//
// It returns a map containing the weather data and an error if any occurred
// during the process.
//
// The function performs the following steps:
//  1. Logs the initiation of the weather data fetching process.
//  2. Constructs the API request URL using the base URL and query parameters.
//  3. Logs the sanitized request URL (without the API key).
//  4. Sends an HTTP GET request to the OpenWeather API.
//  5. Handles any errors that occur during the HTTP request.
//  6. Decodes the JSON response from the API into a map.
//  7. Logs any errors that occur during the JSON decoding process.
//  8. Stores the fetched weather data in a cache.
//  9. Logs the successful fetching and caching of the weather data.
//
// Returns:
// - map[string]interface{}: The weather data fetched from the API.
// - error: An error if any occurred during the process.
func GetWeatherData() (map[string]interface{}, error) {
	utils.LogInfo("Fetching weather data", nil)
	baseURL := "https://api.openweathermap.org/data/2.5/weather"
	queryParams := fmt.Sprintf("q=Copenhagen&appid=%s", config.AppConfig.WeatherAPI.OpenWeatherAPIKey)
	fullURL := fmt.Sprintf("%s?%s", baseURL, queryParams)

	utils.LogInfo("Sending request to OpenWeather API", logrus.Fields{
		"url": fmt.Sprintf("%s?q=Copenhagen", baseURL),
	})

	res, err := http.Get(fullURL)
	if err != nil {
		utils.LogError(err, "Failed to fetch weather data from API", logrus.Fields{"url": baseURL})
		return nil, err
	}
	defer res.Body.Close()

	var weatherData map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&weatherData); err != nil {
		utils.LogError(err, "Failed to decode weather API response", nil)
		return nil, err
	}

	WeatherCache.Set(weatherDataCacheKey, weatherData, weatherDataCacheTime)
	utils.LogInfo("Weather data fetched and stored in cache", nil)

	return weatherData, nil
}