package module

import (
	"context"
	"github.com/Orion777-cmd/weather-app/internal/constants/models"
)

type WeatherService interface {
	GetWeather(ctx context.Context, rq models.WeatherRequest) (models.WeatherResponse, error)
}