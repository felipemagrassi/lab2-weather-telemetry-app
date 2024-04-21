package service

import "testing"

func TestGetTemperature(t *testing.T) {
	service := NewMemoryCepService()
	res, err := service.GetTemperature("00000000")
	if err != CepNotFoundError {
		t.Fatal("invalid error getting temperature for 00000-000")
	}
	if res != nil {
		t.Fatal("result not nil for getting temperature for 00000-000")
	}

	res, err = service.GetTemperature("0")
	if err != InvalidCepError {
		t.Fatal("invalid error getting temperature for 0")
	}
	if res != nil {
		t.Fatal("result not nil for getting temperature for 0")
	}

	res, err = service.GetTemperature("20561250")
	if err != nil {
		t.Fatal("invalid error for getting temperature at 20561250")
	}

	if res.Cep != "20561250" {
		t.Fatal("Invalid CEP for 20561250")
	}
	if res.City != "Rio de Janeiro" {
		t.Fatal("Invalid City for 20561250")
	}

	if res.Temp_C <= 20.0 || res.Temp_C >= 21.0 {
		t.Fatal("Invalid Temp_C for 20561250", res.Temp_C)
	}
	if res.Temp_K <= 293 || res.Temp_K >= 294.0 {
		t.Fatal("Invalid Temp_K for 20561250", res.Temp_K)
	}
	if res.Temp_F <= 68.0 || res.Temp_F >= 70.0 {
		t.Fatal("Invalid Temp_F for 20561250", res.Temp_F)
	}

}
