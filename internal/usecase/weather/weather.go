package weather

import (
	"context"
	"github.com/hong195/wheater-bot/internal/entity"
)

// UseCase -.
type UseCase struct {
	cityDetailsRepo entity.CityDetailRepository
	weatherRepo     entity.WeatherRepository
}

// New -.
func New(cr entity.CityDetailRepository, wr entity.WeatherRepository) *UseCase {
	return &UseCase{
		cityDetailsRepo: cr,
		weatherRepo:     wr,
	}
}

func (uc *UseCase) GetWeatherByCoordinates(ctx context.Context, lat, lon string) (*Weather, error) {

	weather, err := uc.weatherRepo.GetWeatherByCoordinates(ctx, lon, lat)

	if err != nil {
		return nil, err
	}

	cityDetails, err := uc.cityDetailsRepo.GetCityDetailsByCoordinates(ctx, lat, lon)

	if err != nil {
		return nil, err
	}

	return &Weather{
		weather.Temperature,
		weather.FeelsLike,
		weather.Humidity,
		weather.WindSpeed,
		weather.Description,
		cityDetails.City,
		cityDetails.Country,
	}, nil
}
