package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/felipemagrassi/lab2-weather-telemetry-app/service-b/internal/handler"
	"github.com/felipemagrassi/lab2-weather-telemetry-app/service-b/internal/service"
	"github.com/felipemagrassi/lab2-weather-telemetry-app/service-b/internal/usecase"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

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
				semconv.ServiceNameKey.String("service-b"),
			),
		),
	)
	otel.SetTracerProvider(tp)

	return tp.Shutdown, nil
}

func main() {
	var (
		webServerPort                = viper.GetString("HTTP_PORT")
		weatherApiKey                = viper.GetString("WEATHER_API_KEY")
		weatherService               = service.NewWeatherApiService(weatherApiKey)
		cepService                   = service.NewViaCepService()
		getTemperatureFromCepUseCase = usecase.NewGetTemperatureFromCepUseCase(cepService, weatherService)
		getTemperatureHandler        = handler.NewGetTemperatureHandler(getTemperatureFromCepUseCase)
	)

	_, err := initProvider()
	if err != nil {
		fmt.Println("Error initializing provider: ", err)
		return
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/", getTemperatureHandler.Handle)

	fmt.Println("Server running at :", webServerPort)
	port := fmt.Sprintf(":%s", webServerPort)
	http.ListenAndServe(port, r)
}
