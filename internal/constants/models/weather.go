package models

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation"
)

type Location struct {
	Latitude float64 `json:"latitude" bson:"latitude"`
	Longitude float64 `json:"longitude" bson:"longitude"`
}

type WeatherRequest struct {
	City string `json:"location" bson:"location"`
	Coordinate Location `json:"location" bson:"location"`
	DateTime string `json:"datetime" bson:"datetime"`
}

type Weather struct {
	Datetime  string    `json:"datetime" bson:"datetime"`
	Tempmin   float32   `json:"tempmin" bson:"tempmin"`
	Tempmax   float32   `json:"tempmax" bson:"tempmax"`
	Humidity  float32   `json:"humidity" bson:"humidity"`
	Precip    float32   `json:"precip" bson:"precip"`
	Snow      float32   `json:"snow" bson:"snow"`
	Snowdepth float32   `json:"snowdepth" bson:"snowdepth"`
	Windspeed float32   `json:"windspeed" bson:"windspeed"`
	Temp      float32   `json:"temp" bson:"temp"`
	Hours     []Weather `json:"hours" bson:"hour"`
}

type WeatherResponse struct {
	Days []Weather `json:"days" bson:"days"`
}

func (w WeatherRequest) Validate() error {
    err := validation.ValidateStruct(&w,
        validation.Field(&w.DateTime, validation.Required.Error("datetime field required")),
    )
    if err != nil {
        return err
    }

    hasCity := w.City != ""
    hasLocation := w.Coordinate.Latitude != 0 || w.Coordinate.Longitude != 0

    if hasCity == hasLocation {
        return errors.New("either city or location coordinates must be provided, but not both")
    }
    return nil
}