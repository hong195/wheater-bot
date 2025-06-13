package weather

import (
	"context"
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
		City:        "Tashkent",
		Country:     "Uzbekistan",
	}

	res, err := useCase.GetWeatherByCoordinates(context.TODO(), latitude, longitude)

	assert.NotNil(t, useCase, "use case cannot be nil")
	assert.Nil(t, err, "error should be nil")
	assert.NotNil(t, res, "result should not be nil")
	assert.Equal(t, expectedRes, res, "result should be equal")
}
