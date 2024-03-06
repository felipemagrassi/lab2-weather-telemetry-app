package main

import (
	"fmt"
	"net/http"

	"github.com/felipemagrassi/weather-cep-api/internal/handler"
	"github.com/felipemagrassi/weather-cep-api/internal/service"
	"github.com/felipemagrassi/weather-cep-api/internal/usecase"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	cep_service := service.NewViaCepService()
	getTemperatureFromCepUseCase := usecase.NewGetTemperatureFromCepUseCase(cep_service)
	getTemperatureHandler := handler.NewGetTemperatureHandler(getTemperatureFromCepUseCase)

	http.HandleFunc("/", getTemperatureHandler.Handle)
	fmt.Println("Server running at 8080")
	http.ListenAndServe(":8080", nil)
}
