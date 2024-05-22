package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"time"

	"github.com/felipemagrassi/lab2-weather-telemetry-app/service-a/internal/service"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
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

func initProvider(serviceName, collectorURL string) (func(context.Context) error, error) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	conn, err := grpc.NewClient(collectorURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tracerProvider.Shutdown, nil
}

func main() {
	var (
		webServerPort = viper.GetString("HTTP_PORT")
		cepService    = cepServiceGateway(viper.GetString("CEP_SERVICE"))
		spanName      = "CEP_SERVICE_A"
		serviceName   = "servicea"
		otlpEndpoint  = viper.GetString("OTEL_EXPORTER_OTLP_ENDPOINT")
	)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	shutdown, err := initProvider(
		serviceName,
		otlpEndpoint,
	)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("Failed to shutdown trace provider: %w", err)
		}
	}()

	tracer := otel.Tracer("cep-tracer")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		carrier := propagation.HeaderCarrier(r.Header)
		ctx := r.Context()
		ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

		fmt.Printf("\nReceived POST request using service: %s", cepService.Name())
		ctx, initialSpan := tracer.Start(
			ctx,
			fmt.Sprintf("Chamada Externa para %s %s", cepService.Name(), spanName),
		)
		defer initialSpan.End()

		input := &Input{}
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if ok := validCEP(input.Cep); ok != true {
			http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
			return
		}

		output, err := cepService.GetTemperature(ctx, input.Cep)
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
	})

	go func() {
		fmt.Println("Server running at port:", webServerPort)
		port := fmt.Sprintf(":%s", webServerPort)
		if err := http.ListenAndServe(port, nil); err != nil {
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

func cepServiceGateway(cepService string) service.CepService {
	switch strings.ToUpper(cepService) {
	case "MEMORY":
		return service.NewMemoryCepService()
	case "B":
		return service.NewBService()
	default:
		log.Fatal("Invalid CEP Service")
		return nil
	}
}
