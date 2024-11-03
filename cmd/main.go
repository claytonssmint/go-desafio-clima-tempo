package main

import (
	"log"
	"net/http"
	"os"

	"github.com/claytonssmint/clima-tempo-go/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := chi.NewRouter()
	r.Get("/weather", handlers.GetWeatherHandler)

	port := os.Getenv("WEB_SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	address := ":" + port
	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(address, r))
}
