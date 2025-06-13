package webapi

import (
	"context"
	"github.com/hong195/wheater-bot/internal/entity"
)

type InMemoryWeatherRepo struct {
	storage map[entity.Coords]*entity.Weather
}

func NewInMemoryWeatherRepo() *InMemoryWeatherRepo {
	storage := getFakeWeatherDetails()
	return &InMemoryWeatherRepo{
		storage: storage,
	}
}

func getFakeWeatherDetails() map[entity.Coords]*entity.Weather {
	storage := make(map[entity.Coords]*entity.Weather)

	storage[entity.Coords{Longitude: 41.311081, Latitude: 69.240562}] = &entity.Weather{
		Sunrise:     1749729167,
		Sunset:      1749729167,
		Temperature: 30,
		FeelsLike:   0,
		Humidity:    0,
		WindSpeed:   0,
		Description: "Cloudy",
	}

	return storage
}

func (i *InMemoryWeatherRepo) GetWeatherByCoordinates(ctx context.Context, lat, lon float64) (*entity.Weather, error) {

	for coordinates, weatherDetails := range i.storage {

		if coordinates.Latitude == lat && coordinates.Longitude == lon {
			return weatherDetails, nil
		}

	}

	return nil, entity.NotFoundError
}
