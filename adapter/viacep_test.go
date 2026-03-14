package adapter

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/caio/weather-api/domain"
	"github.com/stretchr/testify/assert"
)

func TestViaCEPAdapter_FindLocation_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := ViaCEPResponse{
			Cep:        "01310100",
			Localidade: "São Paulo",
			Uf:         "SP",
			Erro:       false,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	adapter := NewViaCEPAdapter()
	adapter.baseURL = server.URL

	location, err := adapter.FindLocation("01310100")

	assert.NoError(t, err)
	assert.NotNil(t, location)
	assert.Equal(t, "São Paulo", location.City)
	assert.Equal(t, "SP", location.State)
}

func TestViaCEPAdapter_FindLocation_InvalidZipCode(t *testing.T) {
	adapter := NewViaCEPAdapter()

	tests := []struct {
		name    string
		zipCode string
	}{
		{"Too short", "1234567"},
		{"Too long", "123456789"},
		{"Empty", ""},
		{"With letters", "0131010A"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			location, err := adapter.FindLocation(tt.zipCode)
			assert.Error(t, err)
			assert.Nil(t, location)
			assert.Equal(t, domain.ErrInvalidZipCode, err)
		})
	}
}

func TestViaCEPAdapter_FindLocation_NotFound(t *testing.T) {
	// Mock server returning error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := ViaCEPResponse{
			Erro: true,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	adapter := NewViaCEPAdapter()
	adapter.baseURL = server.URL

	location, err := adapter.FindLocation("99999999")

	assert.Error(t, err)
	assert.Nil(t, location)
	assert.Equal(t, domain.ErrZipCodeNotFound, err)
}

func TestViaCEPAdapter_FindLocation_NotFoundStringError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"erro":"true"}`))
	}))
	defer server.Close()

	adapter := NewViaCEPAdapter()
	adapter.baseURL = server.URL

	location, err := adapter.FindLocation("99999999")

	assert.Error(t, err)
	assert.Nil(t, location)
	assert.Equal(t, domain.ErrZipCodeNotFound, err)
}

func TestViaCEPAdapter_FindLocation_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	adapter := NewViaCEPAdapter()
	adapter.baseURL = server.URL

	location, err := adapter.FindLocation("01310100")

	assert.Error(t, err)
	assert.Nil(t, location)
	assert.Equal(t, domain.ErrZipCodeNotFound, err)
}

func TestViaCEPAdapter_FindLocation_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	adapter := NewViaCEPAdapter()
	adapter.baseURL = server.URL

	location, err := adapter.FindLocation("01310100")

	assert.Error(t, err)
	assert.Nil(t, location)
}

func TestViaCEPAdapter_FindLocation_NetworkError(t *testing.T) {
	adapter := NewViaCEPAdapter()
	adapter.baseURL = "http://invalid-url-that-does-not-exist-12345.com"

	location, err := adapter.FindLocation("01310100")

	assert.Error(t, err)
	assert.Nil(t, location)
}
