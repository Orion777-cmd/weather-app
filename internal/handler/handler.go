package handler

import (
	"net/http"
	"time"
	"context"
	"strings"
	"strconv"

	"github.com/Orion777-cmd/weather-app/internal/constants/models"
	"github.com/Orion777-cmd/weather-app/internal/module"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// WeatherHandler handles HTTP requests for weather-related endpoints.
type WeatherHandler struct {
	weatherService module.WeatherService
	logger         *zap.Logger
}

// NewWeatherHandler creates a new WeatherHandler.
func NewWeatherHandler(weatherService module.WeatherService, logger *zap.Logger) *WeatherHandler {
	return &WeatherHandler{
		weatherService: weatherService,
		logger:         logger,
	}
}

// GetWeather handles GET /weather requests.
func (h *WeatherHandler) GetWeather(c *gin.Context) {
    city := c.Query("city")
    coordinateStr := c.Query("coordinate")
    datetime := c.Query("datetime")

	h.logger.Info("abiy Received weather request", zap.String("city", city), zap.String("coordinate", coordinateStr), zap.String("datetime", datetime))
    var location models.Location
    if coordinateStr != "" {
        parts := strings.Split(coordinateStr, ",")
        if len(parts) == 2 {
            lat, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
            lon, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
            if err1 == nil && err2 == nil {
                location = models.Location{Latitude: lat, Longitude: lon}
            }
        }
    }

    rq := models.WeatherRequest{
        City:     city,
        Coordinate: location,
        DateTime: datetime,
    }
    if rq.DateTime == "" {
        rq.DateTime = time.Now().Format("2006-01-02")
    }

    weather, err := h.weatherService.GetWeather(c.Request.Context(), rq)
    if err != nil {
        h.logger.Error("Failed to fetch weather", zap.Error(err), zap.Any("request", rq))
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    h.logger.Info("Weather retrieved", zap.Any("request", rq))
    c.JSON(http.StatusOK, weather)
}

// GetHistory handles GET /history requests.
func (h *WeatherHandler) GetHistory(c *gin.Context) {
	// Assuming Persistence is injected or accessible via WeatherService
	// For simplicity, we'll call a method on WeatherService (to be added)
	history, err := h.weatherService.(interface {
		GetWeatherHistory(ctx context.Context) ([]interface{}, error)
	}).GetWeatherHistory(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to retrieve history", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("History retrieved")
	c.JSON(http.StatusOK, history)
}