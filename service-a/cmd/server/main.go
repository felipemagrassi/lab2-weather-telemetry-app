package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/felipemagrassi/weather-cep-api/service-a/internal/service"
	"github.com/joho/godotenv"
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

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var (
		webServerPort = os.Getenv("WEB_SERVER_PORT")
		mockedService = os.Getenv("MOCKED_SERVICE")
		cepService    = cepServiceGateway(mockedService == "true")
	)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		fmt.Printf("Received POST REQUEST using service: ", cepService)

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

	if err := http.ListenAndServe(fmt.Sprintf(":%s", webServerPort), nil); err != nil {
		log.Fatal("Error initializing server")
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

func cepServiceGateway(mocked bool) service.CepService {
	if mocked {
		return service.NewMemoryCepService()
	}

	return service.NewBService()
}
