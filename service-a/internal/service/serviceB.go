package service

import "context"

type BService struct{}

func NewBService() *BService {
	return &BService{}
}

func (b *BService) Name() string {
	return "B Cep Service"
}

func (b *BService) GetTemperature(ctx context.Context, cep string) (*CepServiceOutput, error) {
	return &CepServiceOutput{}, nil
}
