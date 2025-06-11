package v1

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/hong195/wheater-bot/internal/controller/http/v1/request"
	"net/http"
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

	if err := ctx.QueryParser(&body); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Invalid query")
	}

	weatherDetails, err := r.w.GetWeatherByCoordinates(ctx.Context(), body.Lat, body.Lon)

	fmt.Println(weatherDetails)

	if err != nil {
		r.l.Error(err, "http - v1 - weatherByCoordinates")

		return errorResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	return ctx.Status(http.StatusOK).JSON(weatherDetails)
}
