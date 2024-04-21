package service

import "errors"

type CepService interface {
	GetTemperature(cep string) (*CepServiceOutput, error)
}

type CepServiceOutput struct {
	Cep    string
	City   string
	Temp_C float64
	Temp_K float64
	Temp_F float64
}

var (
	CepNotFoundError = errors.New("Cep Not Found")
	InvalidCepError  = errors.New("Invalid Cep")
)
