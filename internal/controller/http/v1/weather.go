package v1

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/hong195/wheater-bot/internal/controller/http/v1/request"
)

// weatherByCoordinates @Summary     Show city weather
// @Description Detects the city and the weather
// @ID          weather
// @Tags  	    weather
// @Accept      json
// @Produce     json
// @Success     200 {object} weather.Weather
// @Failure     500 {object} response.Error
// @Router      /weather [get]
func (r *V1) weatherByCoordinates(ctx *fiber.Ctx) error {
	var body request.Weather

	if err := ctx.BodyParser(&body); err != nil {
		r.l.Error(err, "http - v1 - weatherByCoordinates")

		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	weatherDetails, err := r.w.GetWeatherByCoordinates(ctx.Context(), body.Lat, body.Lon)

	if err != nil {
		r.l.Error(err, "http - v1 - weatherByCoordinates")

		return errorResponse(ctx, http.StatusInternalServerError, "database problems")
	}

	return ctx.Status(http.StatusOK).JSON(weatherDetails)
}
