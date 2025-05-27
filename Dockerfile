FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o weather-api ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/weather-api .

COPY config/config.yaml ./config/config.yaml

EXPOSE 8080

CMD ["./weather-api"]