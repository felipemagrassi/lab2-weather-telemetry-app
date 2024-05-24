package handler

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/felipemagrassi/lab2-weather-telemetry-app/service-b/internal/usecase"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type GetTemperatureHandler struct {
	getTemperatureFromCep *usecase.GetTemperatureFromCepUseCase
}

type GetTemperatureHandlerOutput struct {
	City       string  `json:"city"`
	Celsius    float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
	Kelvin     float64 `json:"temp_K"`
}

func NewGetTemperatureHandler(getTemperatureFromCep *usecase.GetTemperatureFromCepUseCase) *GetTemperatureHandler {
	return &GetTemperatureHandler{getTemperatureFromCep: getTemperatureFromCep}
}

func (h *GetTemperatureHandler) Handle(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

	cep, ok := h.getCep(r)
	if !ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(`invalid zipcode`))
		return
	}

	input := &usecase.GetTemperatureFromCepInput{Cep: cep}
	output, err := h.getTemperatureFromCep.Execute(ctx, input)
	if err != nil {
		if err == usecase.CepNotFoundError {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&GetTemperatureHandlerOutput{
		City:       output.City,
		Celsius:    output.Celsius,
		Fahrenheit: output.Fahrenheit,
		Kelvin:     output.Kelvin,
	})
}

func (h *GetTemperatureHandler) getCep(r *http.Request) (string, bool) {
	cep := r.URL.Query().Get("cep")

	if cep == "" {
		return "", false
	}

	if len(cep) != 8 {
		return "", false
	}

	numRegex := regexp.MustCompile(`[0-9]`)
	if !numRegex.MatchString(cep) {
		return "", false
	}

	return cep, true
}
