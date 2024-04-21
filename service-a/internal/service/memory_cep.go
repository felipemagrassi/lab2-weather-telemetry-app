package service

import (
	"math/rand"
)

type MemoryCepService struct {
}

func NewMemoryCepService() *MemoryCepService {
	return &MemoryCepService{}
}

func (s *MemoryCepService) GetTemperature(cep string) (*CepServiceOutput, error) {
	if len(cep) != 8 {
		return nil, InvalidCepError
	}

	if cep == "00000000" {
		return nil, CepNotFoundError
	}

	temperature := rand.Float64() + 20

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
