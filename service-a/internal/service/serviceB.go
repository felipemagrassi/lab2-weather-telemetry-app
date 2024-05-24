package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type BService struct{}

func NewBService() *BService {
	return &BService{}
}

func (b *BService) Name() string {
	return "B Cep Service"
}

func (b *BService) GetTemperature(ctx context.Context, cep string) (*CepServiceOutput, error) {
	tr := otel.Tracer("a-b-trace")
	ctx, span := tr.Start(ctx, "BService.GetTemperature")
	defer span.End()

	var output *CepServiceOutput

	url := fmt.Sprintf("http://serviceb:8181/?cep=%s", cep)

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		bytes.NewReader(nil),
	)
	if err != nil {
		return nil, err
	}

	defer request.Body.Close()

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(request.Header))
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return nil, CepNotFoundError
	}

	if response.StatusCode == http.StatusBadRequest {
		return nil, InvalidCepError
	}

	err = json.NewDecoder(response.Body).Decode(&output)
	if err != nil {
		return nil, err
	}

	return output, nil
}
