package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/felipemagrassi/lab2-weather-telemetry-app/service-a/internal/service"
	"github.com/spf13/viper"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type Input struct {
	Cep string `json:"cep"`
}

type ServiceBResponse struct {
	City   string `json:"city"`
	Temp_C string `json:"temp_C"`
	Temp_F string `json:"temp_F"`
	Temp_K string `json:"temp_K"`
}

func init() {
	viper.AutomaticEnv()
}

func initProvider() (func(context.Context) error, error) {
	traceExporter, err := zipkin.New("http://zipkin:9411/api/v2/spans")
	if err != nil {
		return nil, err
	}

	otel.SetTextMapPropagator(propagation.TraceContext{})
	batcher := trace.NewBatchSpanProcessor(traceExporter)

	tp := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithSpanProcessor(batcher),
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String("service-a"),
			),
		),
	)
	otel.SetTracerProvider(tp)

	return tp.Shutdown, nil
}

func main() {
	webServerPort := viper.GetString("HTTP_PORT")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
	)
	defer cancel()

	_, err := initProvider()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		r := chi.NewRouter()
		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)
		r.Post("/cep", cepHandler)

		fmt.Println("Server running at port:", webServerPort)
		port := fmt.Sprintf(":%s", webServerPort)
		if err := http.ListenAndServe(port, r); err != nil {
			log.Fatal("Error initializing server, ", err)
		}
	}()

	select {
	case <-sigCh:
		log.Println("Shutting down gracefully, CTRL+C pressed")
	case <-ctx.Done():
		log.Println("Shutting down because operation is finished")
	}

	_, shutdownCancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer shutdownCancel()
}

func cepHandler(w http.ResponseWriter, r *http.Request) {
	cepService := cepServiceGateway(viper.GetString("CEP_SERVICE"))
	carrier := propagation.HeaderCarrier(
		r.Header,
	)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

	tr := otel.Tracer("a-b-trace")
	ctx, span := tr.Start(ctx, "get weather")
	defer span.End()

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parsedCep, err := parseCep(r.Body)
	if err != nil {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	output, err := cepService.GetTemperature(
		ctx,
		parsedCep,
	)
	if err != nil {
		if err == service.InvalidCepError {
			http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
			return
		}

		if err == service.CepNotFoundError {
			http.Error(w, "can not find zipcode", http.StatusNotFound)
			return
		}

		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&ServiceBResponse{
		City:   output.City,
		Temp_C: fmt.Sprintf("%f", output.Temp_C),
		Temp_K: fmt.Sprintf("%f", output.Temp_K),
		Temp_F: fmt.Sprintf("%f", output.Temp_F),
	})
}

func cepServiceGateway(cepService string) service.CepService {
	switch strings.ToUpper(cepService) {
	case "MEMORY":
		return service.NewMemoryCepService()
	default:
		return service.NewBService()
	}
}

func parseCep(body io.ReadCloser) (string, error) {
	input := &Input{}
	err := json.NewDecoder(body).Decode(&input)
	if err != nil {
		return "", err
	}

	if ok := validCEP(input.Cep); ok != true {
		return "", errors.New("Invalid Cep")
	}

	return input.Cep, nil
}

func validCEP(cep string) bool {
	if len(cep) != 8 {
		return false
	}

	ok, err := regexp.MatchString("^[0-9]*$", cep)
	if err != nil || !ok {
		return false
	}

	return true
}
