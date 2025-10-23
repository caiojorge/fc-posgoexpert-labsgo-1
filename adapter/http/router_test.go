package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetupRouter(t *testing.T) {
	mockUseCase := new(MockGetWeatherUseCase)
	handler := NewWeatherHandler(mockUseCase)

	router := SetupRouter(handler)

	assert.NotNil(t, router)

	routes := router.Routes()
	assert.NotEmpty(t, routes)

	found := false
	for _, route := range routes {
		if route.Path == "/weather/:zipcode" && route.Method == "GET" {
			found = true
			break
		}
	}
	assert.True(t, found, "Weather route should be registered")
}
