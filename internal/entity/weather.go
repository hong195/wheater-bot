package entity

import "context"

type WeatherRepository interface {
	GetWeatherByCoordinates(ctx context.Context, lot, lat string) (*Weather, error)
}

type Weather struct {
	Temperature string
	FeelsLike   string
	Humidity    string
	WindSpeed   string
	Description string
}
