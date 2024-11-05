package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/claytonssmint/clima-tempo-go/services"
	"github.com/claytonssmint/clima-tempo-go/utils"
)

type WeatherResponse struct {
	TempC float64 `json:"temp_c"`
	TempF float64 `json:"temp_f"`
	TempK float64 `json:"temp_k"`
}

func GetWeatherHandler(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Query().Get("cep")
	if !utils.IsvalidCep(cep) {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	log.Printf("Fetching city for CEP: %s", cep)
	city, err := services.GetCityByCEP(cep)
	if err != nil {
		log.Printf("Error fetching city for CEP %s: %v", cep, err)
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}

	log.Printf("City found: %s", city)
	tempC, err := services.GetTemperaturesByCity(city)
	if err != nil {
		log.Printf("Error fetching temperature for city %s: %v", city, err)
		http.Error(w, "erro ao obter temperatura da localidade", http.StatusInternalServerError)
		return
	}

	response := WeatherResponse{
		TempC: tempC,
		TempF: utils.ConverToFahrenheit(tempC),
		TempK: utils.ConverToKelvin(tempC),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
