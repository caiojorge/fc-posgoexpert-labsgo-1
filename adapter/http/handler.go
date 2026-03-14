package http

import (
	"errors"
	"log"
	"net/http"
	"regexp"

	"github.com/caio/weather-api/domain"
	"github.com/caio/weather-api/usecase"
	"github.com/gin-gonic/gin"
)

type WeatherHandler struct {
	getWeatherUseCase usecase.GetWeatherUseCase
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func NewWeatherHandler(getWeatherUseCase usecase.GetWeatherUseCase) *WeatherHandler {
	return &WeatherHandler{
		getWeatherUseCase: getWeatherUseCase,
	}
}

// GetWeatherByZipCode obtém a temperatura de uma cidade pelo CEP
// @Summary Obter temperatura por CEP
// @Description Retorna a temperatura atual em Celsius, Fahrenheit e Kelvin para um CEP brasileiro
// @Tags weather
// @Accept json
// @Produce json
// @Param zipcode path string true "CEP brasileiro (8 dígitos)" example(01310-100)
// @Success 200 {object} domain.Temperature "Temperatura atual"
// @Failure 422 {object} ErrorResponse "CEP inválido"
// @Failure 404 {object} ErrorResponse "CEP não encontrado"
// @Failure 500 {object} ErrorResponse "Erro interno do servidor"
// @Router /weather/{zipcode} [get]
func (h *WeatherHandler) GetWeatherByZipCode(c *gin.Context) {
	zipCode := c.Param("zipcode")

	if !isValidZipCode(zipCode) {
		c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
			Message: "invalid zipcode",
		})
		return
	}

	temperature, err := h.getWeatherUseCase.Execute(zipCode)
	if err != nil {
		log.Printf("Error executing use case: %v", err)
		if errors.Is(err, domain.ErrInvalidZipCode) {
			c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
				Message: "invalid zipcode",
			})
			return
		}

		if errors.Is(err, domain.ErrZipCodeNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Message: "can not find zipcode",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, temperature)
}

func isValidZipCode(zipCode string) bool {
	match, _ := regexp.MatchString(`^\d{8}$`, zipCode)
	return match
}
