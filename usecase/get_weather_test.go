package usecase

import (
	"errors"
	"testing"

	"github.com/caio/weather-api/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockZipCodeRepository struct {
	mock.Mock
}

func (m *MockZipCodeRepository) FindLocation(zipCode string) (*domain.Location, error) {
	args := m.Called(zipCode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Location), args.Error(1)
}

type MockWeatherRepository struct {
	mock.Mock
}

func (m *MockWeatherRepository) GetTemperature(city string) (float64, error) {
	args := m.Called(city)
	return args.Get(0).(float64), args.Error(1)
}

func TestGetWeatherByZipCodeUseCase_Execute_Success(t *testing.T) {
	mockZipCodeRepo := new(MockZipCodeRepository)
	mockWeatherRepo := new(MockWeatherRepository)

	mockZipCodeRepo.On("FindLocation", "01310100").Return(&domain.Location{
		City:  "São Paulo",
		State: "SP",
	}, nil)

	mockWeatherRepo.On("GetTemperature", "São Paulo").Return(28.5, nil)

	useCase := NewGetWeatherByZipCodeUseCase(mockZipCodeRepo, mockWeatherRepo)

	result, err := useCase.Execute("01310100")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 28.5, result.Celsius)
	assert.InDelta(t, 83.3, result.Fahrenheit, 0.1)
	assert.Equal(t, 301.5, result.Kelvin)

	mockZipCodeRepo.AssertExpectations(t)
	mockWeatherRepo.AssertExpectations(t)
}

func TestGetWeatherByZipCodeUseCase_Execute_InvalidZipCode(t *testing.T) {
	mockZipCodeRepo := new(MockZipCodeRepository)
	mockWeatherRepo := new(MockWeatherRepository)

	mockZipCodeRepo.On("FindLocation", "1234567").Return(nil, domain.ErrInvalidZipCode)

	useCase := NewGetWeatherByZipCodeUseCase(mockZipCodeRepo, mockWeatherRepo)

	result, err := useCase.Execute("1234567")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrInvalidZipCode, err)

	mockZipCodeRepo.AssertExpectations(t)
}

func TestGetWeatherByZipCodeUseCase_Execute_ZipCodeNotFound(t *testing.T) {
	mockZipCodeRepo := new(MockZipCodeRepository)
	mockWeatherRepo := new(MockWeatherRepository)

	mockZipCodeRepo.On("FindLocation", "99999999").Return(nil, domain.ErrZipCodeNotFound)

	useCase := NewGetWeatherByZipCodeUseCase(mockZipCodeRepo, mockWeatherRepo)

	result, err := useCase.Execute("99999999")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrZipCodeNotFound, err)

	mockZipCodeRepo.AssertExpectations(t)
}

func TestGetWeatherByZipCodeUseCase_Execute_WeatherNotFound(t *testing.T) {
	mockZipCodeRepo := new(MockZipCodeRepository)
	mockWeatherRepo := new(MockWeatherRepository)

	mockZipCodeRepo.On("FindLocation", "01310100").Return(&domain.Location{
		City:  "São Paulo",
		State: "SP",
	}, nil)

	mockWeatherRepo.On("GetTemperature", "São Paulo").Return(0.0, domain.ErrWeatherNotFound)

	useCase := NewGetWeatherByZipCodeUseCase(mockZipCodeRepo, mockWeatherRepo)

	result, err := useCase.Execute("01310100")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrWeatherNotFound, err)

	mockZipCodeRepo.AssertExpectations(t)
	mockWeatherRepo.AssertExpectations(t)
}

func TestGetWeatherByZipCodeUseCase_Execute_ZipCodeRepoError(t *testing.T) {
	mockZipCodeRepo := new(MockZipCodeRepository)
	mockWeatherRepo := new(MockWeatherRepository)

	expectedError := errors.New("network error")
	mockZipCodeRepo.On("FindLocation", "01310100").Return(nil, expectedError)

	useCase := NewGetWeatherByZipCodeUseCase(mockZipCodeRepo, mockWeatherRepo)

	result, err := useCase.Execute("01310100")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)

	mockZipCodeRepo.AssertExpectations(t)
}

func TestGetWeatherByZipCodeUseCase_Execute_WeatherRepoError(t *testing.T) {
	mockZipCodeRepo := new(MockZipCodeRepository)
	mockWeatherRepo := new(MockWeatherRepository)

	mockZipCodeRepo.On("FindLocation", "01310100").Return(&domain.Location{
		City:  "São Paulo",
		State: "SP",
	}, nil)

	expectedError := errors.New("api error")
	mockWeatherRepo.On("GetTemperature", "São Paulo").Return(0.0, expectedError)

	useCase := NewGetWeatherByZipCodeUseCase(mockZipCodeRepo, mockWeatherRepo)

	result, err := useCase.Execute("01310100")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)

	mockZipCodeRepo.AssertExpectations(t)
	mockWeatherRepo.AssertExpectations(t)
}
