package v1

import (
	"github.com/go-playground/validator/v10"
	"github.com/hong195/wheater-bot/internal/usecase/weather"
	"github.com/hong195/wheater-bot/pkg/logger"
)

// V1 -.
type V1 struct {
	w *weather.UseCase
	l logger.Interface
	v *validator.Validate
}
