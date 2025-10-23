package main

import (
	"log"
	"os"

	"github.com/caio/weather-api/adapter"
	httpAdapter "github.com/caio/weather-api/adapter/http"
	_ "github.com/caio/weather-api/docs"
	"github.com/caio/weather-api/usecase"
)

// @title Weather API
// @version 1.0
// @description API para consulta de temperatura por CEP
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@weather-api.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /

func main() {
	weatherAPIKey := os.Getenv("WEATHER_API_KEY")
	if weatherAPIKey == "" {
		log.Fatal("WEATHER_API_KEY environment variable is required")
	}

	viaCEPAdapter := adapter.NewViaCEPAdapter()
	weatherAPIAdapter := adapter.NewWeatherAPIAdapter(weatherAPIKey)

	getWeatherUseCase := usecase.NewGetWeatherByZipCodeUseCase(
		viaCEPAdapter,
		weatherAPIAdapter,
	)

	weatherHandler := httpAdapter.NewWeatherHandler(getWeatherUseCase)

	router := httpAdapter.SetupRouter(weatherHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
