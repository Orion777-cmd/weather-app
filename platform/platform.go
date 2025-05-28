package platform

import (
	"context"

	"github.com/Orion777-cmd/weather-app/internal/constants/models"
)

type WeatherAPI interface {
	GetWeather(ctx context.Context, rq models.WeatherRequest, response *models.WeatherResponse) error
}