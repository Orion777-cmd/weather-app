package module

import (
	"context"

	"github.com/Orion777-cmd/weather-app/internal/constants/models"
	"github.com/Orion777-cmd/weather-app/platform"
	"github.com/Orion777-cmd/weather-app/internal/repository"
	"go.uber.org/zap"
)

type serviceModule struct {
    log        *zap.Logger
    weatherAPI platform.WeatherAPI
	repo       *repository.WeatherRepository
}

func NewService(weatherAPI platform.WeatherAPI, repo *repository.WeatherRepository, log *zap.Logger) WeatherService {
    return &serviceModule{
        log:        log,
        weatherAPI: weatherAPI,
		repo: 	    repo,
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

	if err := s.repo.SaveWeatherQuery(ctx, rq.City, weatherResponse); err != nil {
        s.log.Error("Failed to save weather query", zap.Error(err))
    }

    return weatherResponse, nil
}