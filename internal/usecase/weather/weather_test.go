package weather

import (
	"context"
	"github.com/hong195/wheater-bot/internal/entity"
	"github.com/hong195/wheater-bot/internal/repo/webapi"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	latitude  = 41.311081
	longitude = 69.240562
)

func getUseCase() *UseCase {
	weatherRepo := webapi.NewInMemoryWeatherRepo()
	cityDetailRepo := webapi.NewInMemoryCityRepo()

	return NewUseCase(cityDetailRepo, weatherRepo)
}

func TestGetWeatherDetailsByCoordinates(t *testing.T) {

	useCase := getUseCase()

	expectedRes := &Weather{
		Sunrise:     1749729167,
		Sunset:      1749729167,
		Temperature: 30,
		FeelsLike:   0,
		Humidity:    0,
		WindSpeed:   0,
		Description: "Cloudy",

		City:         "Tashkent",
		Country:      "Uzbekistan",
		County:       "Tashkent",
		Neighborhood: "Tashkent",
	}

	res, err := useCase.GetWeatherByCoordinates(context.TODO(), latitude, longitude)

	assert.NotNil(t, useCase, "use case cannot be nil")
	assert.Nil(t, err, "error should be nil")
	assert.NotNil(t, res, "result should not be nil")
	assert.Equal(t, expectedRes, res, "result should be equal")
}

func TestItFailsIfRepoCallFails(t *testing.T) {

	type testCase struct {
		Name        string
		Coords      entity.Coords
		ExpectedErr error
	}

	testCases := []testCase{
		{
			Name: "invalid coords",
			Coords: entity.Coords{
				Latitude:  0,
				Longitude: 0,
			},
			ExpectedErr: entity.NotFoundError,
		},
		{
			Name: "invalid latitude",
			Coords: entity.Coords{
				Latitude:  0,
				Longitude: longitude,
			},
			ExpectedErr: entity.NotFoundError,
		},
	}

	useCase := getUseCase()
	for _, tc := range testCases {
		res, err := useCase.GetWeatherByCoordinates(context.TODO(), tc.Coords.Latitude, tc.Coords.Longitude)

		assert.Nil(t, res, "result should be nil, one the repo call fails")
		assert.NotNil(t, err, "error should not be nil")
	}
}
