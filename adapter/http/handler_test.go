package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/caio/weather-api/domain"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockGetWeatherUseCase struct {
	mock.Mock
}

func (m *MockGetWeatherUseCase) Execute(zipCode string) (*domain.Temperature, error) {
	args := m.Called(zipCode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Temperature), args.Error(1)
}

func TestWeatherHandler_GetWeatherByZipCode_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := new(MockGetWeatherUseCase)
	expectedTemp := &domain.Temperature{
		Celsius:    28.5,
		Fahrenheit: 83.3,
		Kelvin:     301.5,
	}

	mockUseCase.On("Execute", "01310100").Return(expectedTemp, nil)

	handler := NewWeatherHandler(mockUseCase)
	router := gin.New()
	router.GET("/weather/:zipcode", handler.GetWeatherByZipCode)

	req, _ := http.NewRequest("GET", "/weather/01310100", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response domain.Temperature
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedTemp.Celsius, response.Celsius)
	assert.Equal(t, expectedTemp.Fahrenheit, response.Fahrenheit)
	assert.Equal(t, expectedTemp.Kelvin, response.Kelvin)

	mockUseCase.AssertExpectations(t)
}

func TestWeatherHandler_GetWeatherByZipCode_InvalidFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := new(MockGetWeatherUseCase)
	handler := NewWeatherHandler(mockUseCase)
	router := gin.New()
	router.GET("/weather/:zipcode", handler.GetWeatherByZipCode)

	tests := []struct {
		name    string
		zipcode string
	}{
		{"Too short", "1234567"},
		{"Too long", "123456789"},
		{"With letters", "0131010A"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/weather/"+tt.zipcode, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

			var response ErrorResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "invalid zipcode", response.Message)
		})
	}
}

func TestWeatherHandler_GetWeatherByZipCode_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := new(MockGetWeatherUseCase)
	mockUseCase.On("Execute", "99999999").Return(nil, domain.ErrZipCodeNotFound)

	handler := NewWeatherHandler(mockUseCase)
	router := gin.New()
	router.GET("/weather/:zipcode", handler.GetWeatherByZipCode)

	req, _ := http.NewRequest("GET", "/weather/99999999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "can not find zipcode", response.Message)

	mockUseCase.AssertExpectations(t)
}

func TestWeatherHandler_GetWeatherByZipCode_InvalidZipCodeError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := new(MockGetWeatherUseCase)
	mockUseCase.On("Execute", "12345678").Return(nil, domain.ErrInvalidZipCode)

	handler := NewWeatherHandler(mockUseCase)
	router := gin.New()
	router.GET("/weather/:zipcode", handler.GetWeatherByZipCode)

	req, _ := http.NewRequest("GET", "/weather/12345678", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "invalid zipcode", response.Message)

	mockUseCase.AssertExpectations(t)
}

func TestWeatherHandler_GetWeatherByZipCode_InternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := new(MockGetWeatherUseCase)
	mockUseCase.On("Execute", "01310100").Return(nil, errors.New("internal error"))

	handler := NewWeatherHandler(mockUseCase)
	router := gin.New()
	router.GET("/weather/:zipcode", handler.GetWeatherByZipCode)

	req, _ := http.NewRequest("GET", "/weather/01310100", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "internal server error", response.Message)

	mockUseCase.AssertExpectations(t)
}

func TestIsValidZipCode(t *testing.T) {
	tests := []struct {
		name     string
		zipCode  string
		expected bool
	}{
		{"Valid zipcode", "01310100", true},
		{"Valid zipcode with zeros", "00000000", true},
		{"Too short", "1234567", false},
		{"Too long", "123456789", false},
		{"With letters", "0131010A", false},
		{"With special chars", "01310-100", false},
		{"Empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidZipCode(tt.zipCode)
			assert.Equal(t, tt.expected, result)
		})
	}
}
