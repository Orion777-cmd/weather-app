package repository

import (
    "context"
    "encoding/json"
    "github.com/Orion777-cmd/weather-app/internal/constants/models"
    "github.com/Orion777-cmd/weather-app/internal/db"
)

type WeatherRepository struct {
    q db.Querier
}

func NewWeatherRepository(q db.Querier) *WeatherRepository {
    return &WeatherRepository{q: q}
}

func (r *WeatherRepository) SaveWeatherQuery(ctx context.Context, city string, weather models.WeatherResponse) error {
    data, err := json.Marshal(weather)
    if err != nil {
        return err
    }
    _, err = r.q.InsertWeatherQuery(ctx, db.InsertWeatherQueryParams{
        City:        city,
        Column2: data,
    })
    return err
}