# Weather Cep API

This API is a simple weather API that returns the temperature of a city by receiving a CEP (Brazilian Zip Code) as a parameter and has spans with zipkin.

## Usage

1. Copy `.env.sample` to `.env`
```bash
cp .env.sample .env
```

2. Add your WEATHER_API_KEY to `.env` file
`
WEATHER_API_KEY=XXXXXXXXX
`

3. Run locally with docker-compose
```bash
docker compose up --build
```

4. Add the CEP as a parameter in the URL and the API will return the temperature of the city.
```bash
curl -X POST localhost:8080/cep -d '{"cep": "20561250"}'
```

## Zipkin Traces

Open `localhost:9411` and you should see the traces from your call
[zipkin](screenshots/zipkin.png)

