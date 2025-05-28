package initiator

import (
	md "github.com/Orion777-cmd/weather-app/internal/module"
	"github.com/Orion777-cmd/weather-app/platform"
	"go.uber.org/zap"
)
type WeatherAPI struct {
	GeocodingBaseURL string
	OneCallBaseURL   string
}
// WeatherModule implements the WeatherService interface.
type module struct {
	weatherModule  md.WeatherService
}

// InitWeatherModule initializes the weather module.
func InitWeatherModule(_ Persistence, weatherAPI platform.WeatherAPI, logger *zap.Logger) module {
	logger.Info("Initializing weather module")
	return module{
		weatherModule: md.NewService(weatherAPI, logger),
	}
}