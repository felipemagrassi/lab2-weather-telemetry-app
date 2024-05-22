package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/felipemagrassi/lab2-weather-telemetry-app/service-a/internal/service"
	"github.com/spf13/viper"
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

func main() {
	var (
		webServerPort = viper.GetString("HTTP_PORT")
		cepService    = cepServiceGateway(viper.GetString("CEP_SERVICE"))
	)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		fmt.Printf("\nReceived POST request using service: %s", cepService.Name())

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

		output, err := cepService.GetTemperature(input.Cep)
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

	fmt.Println("Server running at port:", webServerPort)
	port := fmt.Sprintf(":%s", webServerPort)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("Error initializing server, ", err)
	}
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
