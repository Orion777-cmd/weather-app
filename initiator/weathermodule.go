package initiator

import (
	"go.uber.org/zap"
)


type WeatherAPI struct {
	BaseURL string
}

func initWeatherAPI(baseURL string, logger *zap.Logger) *WeatherAPI {
	logger.Info("Initializing Weather API client", zap.String("baseURL", baseURL))
	weatherAPI := &WeatherAPI{BaseURL: baseURL}
	logger.Info("Weather API client initialized")
	return weatherAPI
}