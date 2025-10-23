package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertTemperature(t *testing.T) {
	tests := []struct {
		name     string
		celsius  float64
		expected Temperature
	}{
		{
			name:    "Zero celsius",
			celsius: 0,
			expected: Temperature{
				Celsius:    0,
				Fahrenheit: 32,
				Kelvin:     273,
			},
		},
		{
			name:    "Positive temperature",
			celsius: 25,
			expected: Temperature{
				Celsius:    25,
				Fahrenheit: 77,
				Kelvin:     298,
			},
		},
		{
			name:    "Negative temperature",
			celsius: -10,
			expected: Temperature{
				Celsius:    -10,
				Fahrenheit: 14,
				Kelvin:     263,
			},
		},
		{
			name:    "Decimal temperature",
			celsius: 28.5,
			expected: Temperature{
				Celsius:    28.5,
				Fahrenheit: 83.3,
				Kelvin:     301.5,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertTemperature(tt.celsius)
			assert.Equal(t, tt.expected.Celsius, result.Celsius)
			assert.InDelta(t, tt.expected.Fahrenheit, result.Fahrenheit, 0.1)
			assert.Equal(t, tt.expected.Kelvin, result.Kelvin)
		})
	}
}

func TestErrors(t *testing.T) {
	assert.NotNil(t, ErrInvalidZipCode)
	assert.NotNil(t, ErrZipCodeNotFound)
	assert.NotNil(t, ErrWeatherNotFound)

	assert.Equal(t, "invalid zipcode", ErrInvalidZipCode.Error())
	assert.Equal(t, "can not find zipcode", ErrZipCodeNotFound.Error())
	assert.Equal(t, "weather not found", ErrWeatherNotFound.Error())
}
