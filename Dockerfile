FROM golang:1.22 as build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build cmd/server/main.go

EXPOSE 8080

ENTRYPOINT ["./main"]

