# Weather Cep API

This API is a simple weather API that returns the temperature of a city by receiving a CEP (Brazilian Zip Code) as a parameter.

## Usage

1. Run locally with docker-compose
```bash
docker compose up --build
```

2. Add the CEP as a parameter in the URL and the API will return the temperature of the city.
```bash
curl "http://localhost:8080/?cep=10010-000"
```

## Using in production

Run `cp .env.sample .env` and fill the environment variables with your own values or 
create a `.env` file with an environment variable called `WEATHER_API_KEY` with your Weather API key (https://www.weatherapi.com/)

1. To run the docker compose production file: 

```bash 
docker compose -f docker-compose.prod.yml up
```

## Demo in Google Cloud Run

* With valid cep
https://weather-cep-api-zj3c47ztra-uc.a.run.app/?cep=20561-250

* With invalid cep
https://weather-cep-api-zj3c47ztra-uc.a.run.app/?cep=

* With not found cep 
https://weather-cep-api-zj3c47ztra-uc.a.run.app/?cep=11111-111
