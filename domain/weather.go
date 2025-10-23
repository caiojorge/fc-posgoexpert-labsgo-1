package domain

import "errors"

var (
	ErrInvalidZipCode  = errors.New("invalid zipcode")
	ErrZipCodeNotFound = errors.New("can not find zipcode")
	ErrWeatherNotFound = errors.New("weather not found")
)

type Temperature struct {
	Celsius    float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
	Kelvin     float64 `json:"temp_K"`
}

type Location struct {
	City  string
	State string
}

type ZipCodeRepository interface {
	FindLocation(zipCode string) (*Location, error)
}

type WeatherRepository interface {
	GetTemperature(city string) (float64, error)
}

func ConvertTemperature(celsius float64) Temperature {
	return Temperature{
		Celsius:    celsius,
		Fahrenheit: celsius*1.8 + 32,
		Kelvin:     celsius + 273,
	}
}
