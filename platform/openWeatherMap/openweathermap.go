package openWeatherMap

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/Orion777-cmd/weather-app/internal/constants/models"
	"github.com/Orion777-cmd/weather-app/platform"
	"go.uber.org/zap"
)

// openWeatherMap implements the WeatherAPI interface for OpenWeatherMap.
type OpenWeatherMap struct {
	geocodingBaseURL string
	oneCallBaseURL   string
	log              *zap.Logger
}

// InitOpenWeatherMap initializes the OpenWeatherMap client.
func InitOpenWeatherMap(geocodingBaseURL, oneCallBaseURL string, log *zap.Logger) platform.WeatherAPI {
	return &OpenWeatherMap{
		geocodingBaseURL: geocodingBaseURL,
		oneCallBaseURL:   oneCallBaseURL,
		log:              log,
	}
}

// geocodeResponse holds the Geocoding API response structure.
type geocodeResponse struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// oneCallResponse holds the One Call API response structure.
type oneCallResponse struct {
	Current struct {
		Dt        int64   `json:"dt"`        // Unix timestamp
		Temp      float64 `json:"temp"`      // Kelvin
		Humidity  float64 `json:"humidity"`  // Percentage
		WindSpeed float64 `json:"wind_speed"` // m/s
		Rain      struct {
			OneHour float64 `json:"1h"` // mm
		} `json:"rain"` // Optional
		Snow      struct {
			OneHour float64 `json:"1h"` // mm
		} `json:"snow"` // Optional
	} `json:"current"`
	Hourly []struct {
		Dt        int64   `json:"dt"`
		Temp      float64 `json:"temp"`
		Humidity  float64 `json:"humidity"`
		WindSpeed float64 `json:"wind_speed"`
		Rain      struct {
			OneHour float64 `json:"1h"`
		} `json:"rain"` // Optional
		Snow      struct {
			OneHour float64 `json:"1h"`
		} `json:"snow"` // Optional
	} `json:"hourly"`
	Daily []struct {
		Dt   int64 `json:"dt"`
		Temp struct {
			Min float64 `json:"min"` // Kelvin
			Max float64 `json:"max"` // Kelvin
			Day float64 `json:"day"` // Kelvin
		} `json:"temp"`
		Humidity  float64 `json:"humidity"`
		WindSpeed float64 `json:"wind_speed"`
		Rain      float64 `json:"rain"` // mm, optional
		Snow      float64 `json:"snow"` // mm, optional
	} `json:"daily"`
}

func (o *OpenWeatherMap) GetWeather(ctx context.Context, rq models.WeatherRequest, response *models.WeatherResponse) error {
	// Validate request
	if err := rq.Validate(); err != nil {
		o.log.Error("Invalid request", zap.Error(err), zap.Any("request", rq))
		return fmt.Errorf("validation failed: %v", err)
	}

	var lat, lon float64

	// Determine coordinates
	if rq.City != "" {
		// City-based: Call Geocoding API
		cityQuery := url.QueryEscape(rq.City)
		o.log.Info("*******************Geocoding base URL template", zap.String("geocodingBaseURL", o.geocodingBaseURL))
		url := fmt.Sprintf(o.geocodingBaseURL, cityQuery)
		o.log.Info("Calling Geocoding API", zap.String("url", url))

		resp, err := http.Get(url)
		if err != nil {
			o.log.Error("Unable to get geocoding data", zap.Error(err), zap.Any("request", rq))
			return fmt.Errorf("geocoding request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			o.log.Error("Unexpected geocoding status", zap.Int("status", resp.StatusCode), zap.Any("request", rq))
			return fmt.Errorf("unexpected geocoding status: %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			o.log.Error("Error reading geocoding response", zap.Error(err), zap.Any("request", rq))
			return fmt.Errorf("error reading geocoding response: %v", err)
		}

		var geocodes []geocodeResponse
		if err := json.Unmarshal(body, &geocodes); err != nil {
			o.log.Error("Error unmarshaling geocoding JSON", zap.Error(err), zap.Any("request", rq))
			return fmt.Errorf("error unmarshaling geocoding JSON: %v", err)
		}

		if len(geocodes) == 0 {
			o.log.Error("No geocoding results found", zap.Any("request", rq))
			return fmt.Errorf("no geocoding results for city: %s", rq.City)
		}

		lat = geocodes[0].Lat
		lon = geocodes[0].Lon
	} else {
		// Lat/lon-based
		lat = rq.Coordinate.Latitude
		lon = rq.Coordinate.Longitude
	}

	// Call One Call API
	url := fmt.Sprintf(o.oneCallBaseURL, fmt.Sprintf("%f", lat), fmt.Sprintf("%f", lon))
	o.log.Info("Calling One Call API", zap.String("url", url))

	resp, err := http.Get(url)
	if err != nil {
		o.log.Error("Unable to get weather data", zap.Error(err), zap.Any("request", rq))
		return fmt.Errorf("weather request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		o.log.Error("Unexpected weather status", zap.Int("status", resp.StatusCode), zap.Any("request", rq))
		return fmt.Errorf("unexpected weather status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		o.log.Error("Error reading weather response", zap.Error(err), zap.Any("request", rq))
		return fmt.Errorf("error reading weather response: %v", err)
	}

	var weatherData oneCallResponse
	if err := json.Unmarshal(body, &weatherData); err != nil {
		o.log.Error("Error unmarshaling weather JSON", zap.Error(err), zap.Any("request", rq))
		return fmt.Errorf("error unmarshaling weather JSON: %v", err)
	}

	// Map to WeatherResponse
	response.Days = make([]models.Weather, 0, len(weatherData.Daily)+1)

	// Add current weather as first day
	currentWeather := models.Weather{
		Datetime:  time.Unix(weatherData.Current.Dt, 0).Format("2006-01-02 15:04:05"),
		Tempmin:   float32(weatherData.Current.Temp - 273.15), // Celsius
		Tempmax:   float32(weatherData.Current.Temp - 273.15), // Approximate
		Humidity:  float32(weatherData.Current.Humidity),
		Precip:    float32(weatherData.Current.Rain.OneHour),
		Snow:      float32(weatherData.Current.Snow.OneHour),
		Snowdepth: 0, // Not provided
		Windspeed: float32(weatherData.Current.WindSpeed),
		Temp:      float32(weatherData.Current.Temp - 273.15),
		Hours:     []models.Weather{},
	}
	response.Days = append(response.Days, currentWeather)

	// Add daily forecasts
	for _, daily := range weatherData.Daily {
		dailyWeather := models.Weather{
			Datetime:  time.Unix(daily.Dt, 0).Format("2006-01-02 15:04:05"),
			Tempmin:   float32(daily.Temp.Min - 273.15),
			Tempmax:   float32(daily.Temp.Max - 273.15),
			Humidity:  float32(daily.Humidity),
			Precip:    float32(daily.Rain),
			Snow:      float32(daily.Snow),
			Snowdepth: 0,
			Windspeed: float32(daily.WindSpeed),
			Temp:      float32(daily.Temp.Day - 273.15),
			Hours:     []models.Weather{},
		}
		response.Days = append(response.Days, dailyWeather)
	}

	// Populate Hours for the first day from hourly data
	if len(weatherData.Hourly) > 0 {
		currentWeather.Hours = make([]models.Weather, 0, len(weatherData.Hourly))
		for _, hourly := range weatherData.Hourly[:24] { // Limit to 24 hours
			hourlyWeather := models.Weather{
				Datetime:  time.Unix(hourly.Dt, 0).Format("2006-01-02 15:04:05"),
				Tempmin:   float32(hourly.Temp - 273.15),
				Tempmax:   float32(hourly.Temp - 273.15),
				Humidity:  float32(hourly.Humidity),
				Precip:    float32(hourly.Rain.OneHour),
				Snow:      float32(hourly.Snow.OneHour),
				Snowdepth: 0,
				Windspeed: float32(hourly.WindSpeed),
				Temp:      float32(hourly.Temp - 273.15),
				Hours:     []models.Weather{},
			}
			currentWeather.Hours = append(currentWeather.Hours, hourlyWeather)
		}
		response.Days[0] = currentWeather
	}

	return nil
}