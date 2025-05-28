FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o weather-app ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/weather-app .

COPY config/config.yaml ./config/config.yaml

EXPOSE 8080

CMD ["./weather-app"]