package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

type ViaCEPResponse struct {
	Localidade string `json:"localidade"`
}

type WeatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

func GetCityByCEP(cep string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", errors.New("CEP not found")
	}

	var data ViaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	return data.Localidade, nil
}

func GetTemperaturesByCity(city string) (float64, error) {
	apiKey := os.Getenv("WEATHER_API_KEY")
	resp, err := http.Get(fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, city))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var data WeatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}

	return data.Current.TempC, nil
}
