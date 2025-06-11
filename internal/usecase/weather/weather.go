package weather

import (
	"context"
	"github.com/hong195/wheater-bot/internal/entity"
	"sync"
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

func (uc *UseCase) GetWeatherByCoordinates(ctx context.Context, lat, lon float64) (*Weather, error) {

	var (
		weather        *entity.Weather
		weatherErr     error
		cityDetails    *entity.CityDetails
		cityDetailsErr error
	)

	wg := &sync.WaitGroup{}

	wg.Add(2)
	go func() {
		defer wg.Done()
		weather, weatherErr = uc.weatherRepo.GetWeatherByCoordinates(ctx, lon, lat)
	}()

	go func() {
		defer wg.Done()

		defer func() {
			if r := recover(); r != nil {
				cityDetailsErr = r.(error)
			}
		}()
		cityDetails, cityDetailsErr = uc.cityDetailsRepo.GetCityDetailsByCoordinates(ctx, lat, lon)
	}()

	wg.Wait()

	if weatherErr != nil {
		return nil, weatherErr
	}

	if cityDetailsErr != nil {
		return nil, cityDetailsErr
	}

	return &Weather{
		weather.Sunrise,
		weather.Sunset,
		weather.Temperature,
		weather.FeelsLike,
		weather.Humidity,
		weather.WindSpeed,
		"Cloudy",
		cityDetails.City,
		cityDetails.Country,
	}, nil
}
