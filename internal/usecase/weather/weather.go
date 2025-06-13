package weather

import (
	"context"
	"github.com/hong195/wheater-bot/internal/entity"
	"golang.org/x/sync/errgroup"
)

// UseCase -.
type UseCase struct {
	cityDetailsRepo entity.CityDetailRepository
	weatherRepo     entity.WeatherRepository
}

// NewUseCase -.
func NewUseCase(cr entity.CityDetailRepository, wr entity.WeatherRepository) *UseCase {
	return &UseCase{
		cityDetailsRepo: cr,
		weatherRepo:     wr,
	}
}

func (uc *UseCase) GetWeatherByCoordinates(ctx context.Context, lat, lon float64) (*Weather, error) {
	var (
		weather     *entity.Weather
		cityDetails *entity.CityDetails
	)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		res, err := uc.weatherRepo.GetWeatherByCoordinates(ctx, lon, lat)

		if err != nil {
			return err
		}

		weather = res
		return nil
	})

	g.Go(func() error {
		res, err := uc.cityDetailsRepo.GetCityDetailsByCoordinates(ctx, lon, lat)

		if err != nil {
			return err
		}

		cityDetails = res
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return &Weather{
		weather.Sunrise,
		weather.Sunset,
		weather.Temperature,
		weather.FeelsLike,
		weather.Humidity,
		weather.WindSpeed,
		weather.Description,
		cityDetails.City,
		cityDetails.Country,
	}, nil
}
