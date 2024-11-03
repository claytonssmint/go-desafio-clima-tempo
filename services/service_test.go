package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type MockViaCEPResponse struct {
	Localidade string `json:"localidade"`
}

func TestGetCityByCEP(t *testing.T) {

	mockServer := func(statusCode int, responseBody interface{}) *httptest.Server {
		return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(statusCode)
			json.NewEncoder(w).Encode(responseBody)
		}))
	}

	t.Run("CEP Válido - Retorna Cidade", func(t *testing.T) {
		server := mockServer(http.StatusOK, MockViaCEPResponse{Localidade: "São Paulo"})
		defer server.Close()

		// Passa o URL do mock server como base para a função
		city, err := GetCityByCEP("01001000")

		if err != nil {
			t.Fatalf("Erro ao buscar cidade: %v", err)
		}

		expectedCity := "São Paulo"
		if city != expectedCity {
			t.Errorf("Cidade esperada '%s', mas obteve '%s'", expectedCity, city)
		}
	})

	t.Run("CEP Inválido - Retorna Erro", func(t *testing.T) {
		server := mockServer(http.StatusNotFound, nil)
		defer server.Close()

		_, err := GetCityByCEP("1124")

		if err == nil {
			t.Fatalf("Esperado erro ao buscar cidade, mas não obteve")
		}
	})
}

func TestGetTemperaturesByCity(t *testing.T) {
	mockServer := func(statusCode int, responseBody interface{}) *httptest.Server {
		return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(statusCode)
			json.NewEncoder(w).Encode(responseBody)
		}))
	}

	t.Run("Cidade Inválida - Retorna Erro", func(t *testing.T) {
		server := mockServer(http.StatusNotFound, nil)
		defer server.Close()

		_, err := GetTemperaturesByCity("São Paulo")

		if err == nil {
			t.Fatalf("Esperado erro ao buscar temperatura, mas não obteve")
		}
	})

	t.Run("Cidade Válida - Retorna Temperatura", func(t *testing.T) {
		mockServer2 := func(statusCode int, responseBody interface{}) *httptest.Server {
			return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(statusCode)
				if responseBody != nil {
					json.NewEncoder(w).Encode(responseBody)
				} else {
					w.Write([]byte("{}"))
				}
			}))
		}

		mockResponse := WeatherAPIResponse{
			Current: struct {
				TempC float64 `json:"temp_c"`
			}{
				TempC: 0.00,
			},
		}
		server := mockServer2(http.StatusOK, mockResponse)
		defer server.Close()

		os.Setenv("WEATHER_API_KEY", "mock-api-key")

		apiURL := fmt.Sprintf("%s?key=%s", server.URL, "mock-api-key")

		temperature, err := GetTemperaturesByCity(apiURL)

		if err != nil {
			t.Fatalf("Erro ao buscar temperatura: %v", err)
		}

		expectedTemperature := 0.00
		if temperature != expectedTemperature {
			t.Errorf("Temperatura esperada %.2f, mas obteve %.2f", expectedTemperature, temperature)
		}
	})

}
