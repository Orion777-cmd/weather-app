package module

import (
	"context"

	"github.com/Orion777-cmd/weather-app/internal/constants/models"
	"github.com/Orion777-cmd/weather-app/platform"
	"go.uber.org/zap"
)

type serviceModule struct {
    log        *zap.Logger
    weatherAPI platform.WeatherAPI
}

func NewService(weatherAPI platform.WeatherAPI, log *zap.Logger) WeatherService {
    return &serviceModule{
        log:        log,
        weatherAPI: weatherAPI,
    }
}

func (s *serviceModule) GetWeather(ctx context.Context, rq models.WeatherRequest) (models.WeatherResponse, error) {
    if err := rq.Validate(); err != nil {
        s.log.Warn(err.Error(), zap.Any("request", rq))
        return models.WeatherResponse{}, err
    }

    var weatherResponse models.WeatherResponse
    if err := s.weatherAPI.GetWeather(ctx, rq, &weatherResponse); err != nil {
        return models.WeatherResponse{}, err
    }

    return weatherResponse, nil
}