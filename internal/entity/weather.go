package entity

import "context"

type (
	WeatherRepository interface {
		GetWeatherByCoordinates(ctx context.Context, lat, lot float64) (*Weather, error)
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
