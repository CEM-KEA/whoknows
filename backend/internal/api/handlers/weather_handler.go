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
		http.Error(w, fetchDataError, http.StatusInternalServerError)
		return
	}

	response := WeatherResponse{Data: data}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		utils.LogError(err, encodeResponseError, nil)
		http.Error(w, encodeResponseError, http.StatusInternalServerError)
		return
	}

	utils.LogInfo("Weather data fetched and response sent successfully", nil)
}

// GetWeatherData fetches current weather data for Copenhagen from OpenWeatherMap API
func GetWeatherData() (map[string]interface{}, error) {
	utils.LogInfo("Fetching weather data", nil)
	if cachedData, found := WeatherCache.Get(weatherDataCacheKey); found {
		utils.LogInfo("Weather data found in cache", nil)
		return cachedData.(map[string]interface{}), nil
	}

	apiKey := config.AppConfig.WeatherAPI.OpenWeatherAPIKey
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=Copenhagen&appid=%s", apiKey)
	res, err := http.Get(url)
	if err != nil {
		utils.LogError(err, "Failed to fetch weather data from API", logrus.Fields{
			"url": url,
		})
		return nil, err
	}
	defer res.Body.Close()

	var weatherData map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&weatherData); err != nil {
		utils.LogError(err, "Failed to decode weather API response", nil)
		return nil, err
	}
    
	WeatherCache.Set(weatherDataCacheKey, weatherData, weatherDataCacheTime)
	utils.LogInfo("Weather data fetched and stored in cache successfully", nil)

	return weatherData, nil
}