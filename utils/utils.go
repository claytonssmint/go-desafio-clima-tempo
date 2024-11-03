package utils

import "regexp"

func IsvalidCep(cep string) bool {
	re := regexp.MustCompile(`^\d{8}$`)
	return re.MatchString(cep)
}

func ConverToFahrenheit(celsius float64) float64 {
	return celsius*1.8 + 32
}

func ConverToKelvin(celsius float64) float64 {
	return celsius + 273.15
}
