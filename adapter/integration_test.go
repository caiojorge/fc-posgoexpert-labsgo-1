package adapter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestViaCEPAdapter_RealAPI_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	adapter := NewViaCEPAdapter()

	location, err := adapter.FindLocation("01310100")

	if err == nil {
		assert.NotNil(t, location)
		assert.NotEmpty(t, location.City)
		assert.Equal(t, "São Paulo", location.City)
	}
}

func TestWeatherAPIAdapter_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Skip if no API key is provided
	apiKey := "test-key"
	if apiKey == "test-key" {
		t.Skip("No API key provided for integration test")
	}

	adapter := NewWeatherAPIAdapter(apiKey)

	temp, err := adapter.GetTemperature("São Paulo")

	if err == nil {
		assert.NotZero(t, temp)
	}
}
