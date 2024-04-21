package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/felipemagrassi/lab2-weather-telemetry-app/service-b/internal/handler"
	"github.com/felipemagrassi/lab2-weather-telemetry-app/service-b/internal/service"
	"github.com/felipemagrassi/lab2-weather-telemetry-app/service-b/internal/usecase"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	var (
		webServerPort                = os.Getenv("WEB_SERVER_PORT")
		weatherApiKey                = os.Getenv("WEATHER_API_KEY")
		weatherService               = service.NewWeatherApiService(weatherApiKey)
		cepService                   = service.NewViaCepService()
		getTemperatureFromCepUseCase = usecase.NewGetTemperatureFromCepUseCase(cepService, weatherService)
		getTemperatureHandler        = handler.NewGetTemperatureHandler(getTemperatureFromCepUseCase)
	)

	http.HandleFunc("/", getTemperatureHandler.Handle)
	fmt.Println("Server running at :", webServerPort)
	http.ListenAndServe(fmt.Sprintf(":%s", webServerPort), nil)
}
