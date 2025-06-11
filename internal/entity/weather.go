package entity

import "context"

type (
	WeatherRepository interface {
		GetWeatherByCoordinates(ctx context.Context, lot, lat float64) (*Weather, error)
	}

	Weather struct {
		Sunrise     int
		Sunset      int
		Temperature float64
		FeelsLike   float64
		Humidity    int
		WindSpeed   float64
		Description string
	}
)
