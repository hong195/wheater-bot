package v1

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/hong195/wheater-bot/internal/usecase/weather"
	"github.com/hong195/wheater-bot/pkg/logger"
)

// NewWeatherRoutes -.
func NewWeatherRoutes(apiV1Group fiber.Router, w *weather.UseCase, l logger.Interface) {
	r := &V1{w: w, l: l, v: validator.New(validator.WithRequiredStructEnabled())}

	weatherGroup := apiV1Group.Group("/weather")
	{
		weatherGroup.Get("/", r.weatherByCoordinates)
	}
}
