package usecase

import (
	"context"

	"github.com/felipemagrassi/lab2-weather-telemetry-app/service-b/internal/service"
	"go.opentelemetry.io/otel"
)

type GetTemperatureFromCepInput struct {
	Cep string
}

type GetTemperatureFromCepOutput struct {
	Celsius    float64
	Fahrenheit float64
	Kelvin     float64
	City       string
}

type GetTemperatureFromCepUseCase struct {
	CepService     service.CepService
	WeatherService service.WeatherService
}

func NewGetTemperatureFromCepUseCase(
	cepService service.CepService,
	weatherService service.WeatherService,
) *GetTemperatureFromCepUseCase {
	return &GetTemperatureFromCepUseCase{
		CepService:     cepService,
		WeatherService: weatherService,
	}
}

var CepNotFoundError = service.CepNotFoundError

func (u *GetTemperatureFromCepUseCase) Execute(
	ctx context.Context,
	input *GetTemperatureFromCepInput,
) (*GetTemperatureFromCepOutput, error) {
	tracer := otel.Tracer("a-b-trace")
	ctx, span := tracer.Start(ctx, "GetTemperatureFromCepUseCase.Execute")
	defer span.End()

	address, err := u.CepService.GetAddressByCep(ctx, input.Cep)
	if err != nil {
		return nil, err
	}
	weather, err := u.WeatherService.GetWeatherByCity(ctx, address.Localidade)
	if err != nil {
		return nil, err
	}
	return &GetTemperatureFromCepOutput{
		Celsius:    weather.Temp_c,
		Fahrenheit: weather.Temp_f,
		Kelvin:     weather.Temp_c + 273.15,
		City:       address.Localidade,
	}, nil
}
