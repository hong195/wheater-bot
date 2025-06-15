package webapi

import (
	"context"
	"github.com/hong195/wheater-bot/internal/entity"
)

type InMemoryCityRepo struct {
	storage map[entity.Coords]entity.CityDetails
}

func NewInMemoryCityRepo() *InMemoryCityRepo {
	store := getFakeData()
	return &InMemoryCityRepo{
		store,
	}
}

func getFakeData() map[entity.Coords]entity.CityDetails {
	storage := make(map[entity.Coords]entity.CityDetails)

	storage[entity.Coords{Latitude: 41.311081, Longitude: 69.240562}] = entity.CityDetails{
		City:    "Tashkent",
		Country: "Uzbekistan",
	}

	return storage
}

func (i *InMemoryCityRepo) GetCityDetailsByCoordinates(ctx context.Context, lat, lon float64) (*entity.CityDetails, error) {

	for coordinates, cityDetails := range i.storage {
		if coordinates.Latitude == lat && coordinates.Longitude == lon {
			return &cityDetails, nil
		}

	}

	return nil, entity.NotFoundError
}
