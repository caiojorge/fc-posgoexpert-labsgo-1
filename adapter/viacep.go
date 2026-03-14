package adapter

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/caio/weather-api/domain"
)

type ViaCEPErrorFlag bool

func (f *ViaCEPErrorFlag) UnmarshalJSON(data []byte) error {
	value := strings.TrimSpace(string(data))

	switch value {
	case "true", `"true"`:
		*f = true
		return nil
	case "false", `"false"`, "null", "":
		*f = false
		return nil
	default:
		return fmt.Errorf("invalid ViaCEP erro value: %s", value)
	}
}

type ViaCEPResponse struct {
	Cep         string          `json:"cep"`
	Logradouro  string          `json:"logradouro"`
	Complemento string          `json:"complemento"`
	Bairro      string          `json:"bairro"`
	Localidade  string          `json:"localidade"`
	Uf          string          `json:"uf"`
	Erro        ViaCEPErrorFlag `json:"erro"`
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
		log.Printf("Invalid zipcode: %s", zipCode)
		return nil, domain.ErrInvalidZipCode
	}

	url := fmt.Sprintf("%s/%s/json/", v.baseURL, zipCode)

	resp, err := v.httpClient.Get(url)
	if err != nil {
		log.Printf("Error fetching zipcode from ViaCEP: %v", err)
		return nil, fmt.Errorf("failed to fetch zipcode: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error fetching zipcode from ViaCEP: %v", err)
		return nil, domain.ErrZipCodeNotFound
	}

	var viaCEPResp ViaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&viaCEPResp); err != nil {
		log.Printf("Error decoding response from ViaCEP: %v", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if bool(viaCEPResp.Erro) {
		log.Printf("Zipcode not found: %s", zipCode)
		return nil, domain.ErrZipCodeNotFound
	}

	return &domain.Location{
		City:  viaCEPResp.Localidade,
		State: viaCEPResp.Uf,
	}, nil
}
