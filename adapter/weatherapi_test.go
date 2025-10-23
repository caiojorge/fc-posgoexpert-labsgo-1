package adapter

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/caio/weather-api/domain"
	"github.com/stretchr/testify/assert"
)

func TestWeatherAPIAdapter_GetTemperature_Success(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := WeatherAPIResponse{
			Location: struct {
				Name    string `json:"name"`
				Region  string `json:"region"`
				Country string `json:"country"`
			}{
				Name:    "São Paulo",
				Region:  "Sao Paulo",
				Country: "Brazil",
			},
			Current: struct {
				TempC float64 `json:"temp_c"`
				TempF float64 `json:"temp_f"`
			}{
				TempC: 28.5,
				TempF: 83.3,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	adapter := NewWeatherAPIAdapter("test-key")
	adapter.baseURL = server.URL

	temp, err := adapter.GetTemperature("São Paulo")

	assert.NoError(t, err)
	assert.Equal(t, 28.5, temp)
}

func TestWeatherAPIAdapter_GetTemperature_NotFound(t *testing.T) {
	// Mock server returning 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	adapter := NewWeatherAPIAdapter("test-key")
	adapter.baseURL = server.URL

	temp, err := adapter.GetTemperature("InvalidCity")

	assert.Error(t, err)
	assert.Equal(t, 0.0, temp)
	assert.Equal(t, domain.ErrWeatherNotFound, err)
}

func TestWeatherAPIAdapter_GetTemperature_InvalidJSON(t *testing.T) {
	// Mock server returning invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	adapter := NewWeatherAPIAdapter("test-key")
	adapter.baseURL = server.URL

	temp, err := adapter.GetTemperature("São Paulo")

	assert.Error(t, err)
	assert.Equal(t, 0.0, temp)
}

func TestWeatherAPIAdapter_GetTemperature_ServerError(t *testing.T) {
	// Mock server returning 500
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	adapter := NewWeatherAPIAdapter("test-key")
	adapter.baseURL = server.URL

	temp, err := adapter.GetTemperature("São Paulo")

	assert.Error(t, err)
	assert.Equal(t, 0.0, temp)
	assert.Equal(t, domain.ErrWeatherNotFound, err)
}

func TestWeatherAPIAdapter_GetTemperature_NetworkError(t *testing.T) {
	adapter := NewWeatherAPIAdapter("test-key")
	// Use an invalid URL to trigger network error
	adapter.baseURL = "http://invalid-url-that-does-not-exist-12345.com"

	temp, err := adapter.GetTemperature("São Paulo")

	assert.Error(t, err)
	assert.Equal(t, 0.0, temp)
}
