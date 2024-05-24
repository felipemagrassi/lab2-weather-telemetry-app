package service

import (
	"context"
	"math"
	"math/rand"

	"go.opentelemetry.io/otel"
)

type MemoryCepService struct{}

func NewMemoryCepService() *MemoryCepService {
	return &MemoryCepService{}
}

func (s *MemoryCepService) Name() string {
	return "Memory Cep Service"
}

func (s *MemoryCepService) GetTemperature(ctx context.Context, cep string) (*CepServiceOutput, error) {
	tr := otel.Tracer("a-b-trace")
	ctx, span := tr.Start(ctx, "MemoryCepService.GetTemperature")
	defer span.End()

	if len(cep) != 8 {
		return nil, InvalidCepError
	}

	if cep == "00000000" {
		return nil, CepNotFoundError
	}

	min := 5.0
	max := 35.0
	random := min + rand.Float64()*(max-min)
	temperature := math.Round(random)

	return &CepServiceOutput{
		Cep:    cep,
		Temp_C: temperature,
		Temp_K: toK(temperature),
		Temp_F: toF(temperature),
		City:   "Rio de Janeiro",
	}, nil
}

func toF(celsius float64) float64 {
	return celsius*1.8 + 32
}

func toK(celsius float64) float64 {
	return celsius + 273
}
