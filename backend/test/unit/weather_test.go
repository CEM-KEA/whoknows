package unit_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CEM-KEA/whoknows/backend/internal/api/handlers"
	"github.com/stretchr/testify/assert"
)

var mockConfig = struct {
    AppConfig struct {
        WeatherAPI struct {
            OpenWeatherAPIKey string
        }
    }
}{
    AppConfig: struct {
        WeatherAPI struct {
            OpenWeatherAPIKey string
        }
    }{
        WeatherAPI: struct {
            OpenWeatherAPIKey string
        }{
            OpenWeatherAPIKey: "mock-api-key",
        },
    },
}

func TestGetWeatherData(t *testing.T) {
    // mock openweathermap response
    mockResponse := `{"weather": [{"description": "clear sky"}], "main": {"temp": 280.32}}`
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, mockResponse)
    }))
    defer server.Close()

    // make mockurl to override the original url
    originalURL := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=Copenhagen&appid=%s", mockConfig.AppConfig.WeatherAPI.OpenWeatherAPIKey)
    mockURL := server.URL + "?q=Copenhagen&appid=mock-api-key"

    // mock http.Get
    httpGet = func(url string) (*http.Response, error) {
        if url == originalURL {
            return http.Get(mockURL)
        }
        return http.Get(url)
    }
    defer func() { httpGet = http.Get }() // restore original http.Get

    data, err := handlers.GetWeatherData()
    assert.NoError(t, err)
    assert.NotNil(t, data)

    // verify data is cached
    cachedData, found := handlers.WeatherCache.Get("weatherData")
    assert.True(t, found)
    assert.Equal(t, data, cachedData)

    // call the function again to ensure it returns the cached data
    cachedData, err = handlers.GetWeatherData()
    assert.NoError(t, err)
    assert.Equal(t, data, cachedData)
}

// overrideable http.Get function
var httpGet = http.Get