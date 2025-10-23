package adapter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/caio/weather-api/domain"
)

type ViaCEPResponse struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Erro        bool   `json:"erro"`
}

type ViaCEPAdapter struct {
	baseURL    string
	httpClient *http.Client
}

func NewViaCEPAdapter() *ViaCEPAdapter {
	return &ViaCEPAdapter{
		baseURL: "https://viacep.com.br/ws",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (v *ViaCEPAdapter) FindLocation(zipCode string) (*domain.Location, error) {
	// Validate zipcode: must be exactly 8 digits
	match, _ := regexp.MatchString(`^\d{8}$`, zipCode)
	if !match {
		return nil, domain.ErrInvalidZipCode
	}

	url := fmt.Sprintf("%s/%s/json/", v.baseURL, zipCode)

	resp, err := v.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch zipcode: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, domain.ErrZipCodeNotFound
	}

	var viaCEPResp ViaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&viaCEPResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if viaCEPResp.Erro {
		return nil, domain.ErrZipCodeNotFound
	}

	return &domain.Location{
		City:  viaCEPResp.Localidade,
		State: viaCEPResp.Uf,
	}, nil
}
