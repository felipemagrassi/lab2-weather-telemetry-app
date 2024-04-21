package service

type BService struct {
}

func NewBService() *BService {
	return &BService{}
}

func (b *BService) GetTemperature(cep string) (*CepServiceOutput, error) {
	return &CepServiceOutput{}, nil
}
