package adapter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/caio/weather-api/domain"
)

type WeatherAPIResponse struct {
	Location struct {
		Name    string `json:"name"`
		Region  string `json:"region"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC float64 `json:"temp_c"`
		TempF float64 `json:"temp_f"`
	} `json:"current"`
}

type WeatherAPIAdapter struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

func NewWeatherAPIAdapter(apiKey string) *WeatherAPIAdapter {
	return &WeatherAPIAdapter{
		apiKey:  apiKey,
		baseURL: "https://api.weatherapi.com/v1",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (w *WeatherAPIAdapter) GetTemperature(city string) (float64, error) {
	params := url.Values{}
	params.Add("key", w.apiKey)
	params.Add("q", city)
	params.Add("aqi", "no")

	url := fmt.Sprintf("%s/current.json?%s", w.baseURL, params.Encode())

	resp, err := w.httpClient.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch weather: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, domain.ErrWeatherNotFound
	}

	var weatherResp WeatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherResp); err != nil {
		return 0, fmt.Errorf("failed to decode weather response: %w", err)
	}

	return weatherResp.Current.TempC, nil
}
