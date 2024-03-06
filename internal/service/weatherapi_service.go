package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

type WeatherService interface {
	GetWeatherByCity(ctx context.Context, city string) (*WeatherResponse, error)
}

type WeatherApiService struct {
	apiKey string
	client *http.Client
	logger *log.Logger
}

type WeatherApiResponse struct {
	Location struct {
		Name string `json:"name"`
	} `json:"location"`
	Current struct {
		Temp_c float64 `json:"temp_c"`
		Temp_f float64 `json:"temp_f"`
	}
}

type WeatherResponse struct {
	Name   string  `json:"name"`
	Temp_c float64 `json:"temp_c"`
	Temp_f float64 `json:"temp_f"`
}

func NewWeatherApiService(apiKey string) *WeatherApiService {
	return &WeatherApiService{
		client: &http.Client{},
		apiKey: apiKey,
		logger: log.New(os.Stdout, "weatherapi_service: ", log.LstdFlags),
	}
}

var WeatherServiceError = errors.New("error getting weather")

func (w *WeatherApiService) GetWeatherByCity(ctx context.Context, city string) (*WeatherResponse, error) {
	queryParams := url.Values{}
	queryParams.Add("q", city)
	queryParams.Add("key", w.apiKey)
	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?%s", queryParams.Encode())

	log.Println("Requesting weather data from weatherapi.com")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := w.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, WeatherServiceError
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	weatherApiResponse := &WeatherApiResponse{}
	err = json.Unmarshal(body, &weatherApiResponse)
	if err != nil {
		return nil, err
	}

	weatherResponse := &WeatherResponse{
		Name:   weatherApiResponse.Location.Name,
		Temp_c: weatherApiResponse.Current.Temp_c,
		Temp_f: weatherApiResponse.Current.Temp_f,
	}

	return weatherResponse, nil
}
