package usecase

import "github.com/caio/weather-api/domain"

type GetWeatherUseCase interface {
	Execute(zipCode string) (*domain.Temperature, error)
}
