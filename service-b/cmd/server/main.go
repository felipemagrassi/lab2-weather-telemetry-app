package main

import (
	"fmt"
	"net/http"

	"github.com/felipemagrassi/lab2-weather-telemetry-app/service-b/internal/handler"
	"github.com/felipemagrassi/lab2-weather-telemetry-app/service-b/internal/service"
	"github.com/felipemagrassi/lab2-weather-telemetry-app/service-b/internal/usecase"
	"github.com/spf13/viper"
)

func init() {
	viper.AutomaticEnv()
}

func main() {
	var (
		webServerPort                = viper.GetString("HTTP_PORT")
		weatherApiKey                = viper.GetString("WEATHER_API_KEY")
		weatherService               = service.NewWeatherApiService(weatherApiKey)
		cepService                   = service.NewViaCepService()
		getTemperatureFromCepUseCase = usecase.NewGetTemperatureFromCepUseCase(cepService, weatherService)
		getTemperatureHandler        = handler.NewGetTemperatureHandler(getTemperatureFromCepUseCase)
	)

	http.HandleFunc("/", getTemperatureHandler.Handle)
	fmt.Println("Server running at :", webServerPort)
	port := fmt.Sprintf(":%s", webServerPort)
	http.ListenAndServe(port, nil)
}
