package usecase

import (
	"context"
	"testing"

	"github.com/felipemagrassi/lab2-weather-telemetry-app/service-b/internal/service"
	"github.com/felipemagrassi/lab2-weather-telemetry-app/service-b/internal/service/mocks"
	"go.uber.org/mock/gomock"
)

func TestGetTemperatureFromCepUseCase(t *testing.T) {
	controller := gomock.NewController(t)
	cepService := mocks.NewMockCepService(controller)
	weatherService := mocks.NewMockWeatherService(controller)
	ctx := context.Background()
	cep := "12345678"

	cepService.
		EXPECT().
		GetAddressByCep(ctx, cep).
		Return(&service.ViaCepResponse{
			Cep:         "12345678",
			Logradouro:  "Logradouro",
			Complemento: "Complemento",
			Bairro:      "Bairro",
			Localidade:  "Localidade",
			Uf:          "Uf",
			Ibge:        "Ibge",
			Gia:         "Gia",
			Ddd:         "Ddd",
			Siafi:       "Siafi",
		}, nil)

	weatherService.EXPECT().GetWeatherByCity(ctx, "Localidade").Return(&service.WeatherResponse{
		Name:   "Localidade",
		Temp_c: 10,
		Temp_f: 50,
	}, nil)

	usecase := NewGetTemperatureFromCepUseCase(cepService, weatherService)
	output, err := usecase.Execute(ctx, &GetTemperatureFromCepInput{Cep: cep})
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if output.Celsius != 10 {
		t.Errorf("Expected 10, got %v", output.Celsius)
	}

	if output.Fahrenheit != 50 {
		t.Errorf("Expected 50, got %v", output.Fahrenheit)
	}

	if output.Kelvin != 283.15 {
		t.Errorf("Expected 283.15, got %v", output.Kelvin)
	}

}
