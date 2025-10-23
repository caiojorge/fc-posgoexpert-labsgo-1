package usecase

import (
	"github.com/caio/weather-api/domain"
)

type GetWeatherByZipCodeUseCase struct {
	zipCodeRepo domain.ZipCodeRepository
	weatherRepo domain.WeatherRepository
}

func NewGetWeatherByZipCodeUseCase(
	zipCodeRepo domain.ZipCodeRepository,
	weatherRepo domain.WeatherRepository,
) *GetWeatherByZipCodeUseCase {
	return &GetWeatherByZipCodeUseCase{
		zipCodeRepo: zipCodeRepo,
		weatherRepo: weatherRepo,
	}
}

func (uc *GetWeatherByZipCodeUseCase) Execute(zipCode string) (*domain.Temperature, error) {
	location, err := uc.zipCodeRepo.FindLocation(zipCode)
	if err != nil {
		return nil, err
	}

	celsius, err := uc.weatherRepo.GetTemperature(location.City)
	if err != nil {
		return nil, err
	}

	temperature := domain.ConvertTemperature(celsius)

	return &temperature, nil
}
