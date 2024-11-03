package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
)

type MockServices struct {
	mock.Mock
}

func (m *MockServices) GetCityByCEP(cep string) (string, error) {
	args := m.Called(cep)
	return args.String(0), args.Error(1)
}

func (m *MockServices) GetTemperaturesByCity(city string) (float64, error) {
	args := m.Called(city)
	return args.Get(0).(float64), args.Error(1)
}

func TestGetWeatherHandler(t *testing.T) {
	mockServices := new(MockServices)

	handler := func(w http.ResponseWriter, r *http.Request) {
		GetWeatherHandler(w, r)
	}

	t.Run("CEP Inválido - Retorna Erro", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/weather?cep=invalid", nil)
		if err != nil {
			t.Fatalf("Erro ao criar requisição: %v", err)
		}

		rr := httptest.NewRecorder()
		handler(rr, req)

		if status := rr.Code; status != http.StatusUnprocessableEntity {
			t.Errorf("Esperado status 422, mas obteve %d", status)
		}

		expected := "invalid zipcode\n"
		if rr.Body.String() != expected {
			t.Errorf("Esperado corpo '%s', mas obteve '%s'", expected, rr.Body.String())
		}
	})

	t.Run("Erro ao Obter Temperatura - Retorna Erro", func(t *testing.T) {
		mockServices.On("GetCityByCEP", "01001000").Return("São Paulo", nil)
		mockServices.On("GetTemperaturesByCity", "São Paulo").Return(0.0, errors.New("error"))

		req, err := http.NewRequest(http.MethodGet, "/weather?cep=01001000", nil)
		if err != nil {
			t.Fatalf("Erro ao criar requisição: %v", err)
		}

		rr := httptest.NewRecorder()
		handler(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("Esperado status 500, mas obteve %d", status)
		}

		expected := "erro ao obter temperatura da localidade\n"
		if rr.Body.String() != expected {
			t.Errorf("Esperado corpo '%s', mas obteve '%s'", expected, rr.Body.String())
		}
	})
}
