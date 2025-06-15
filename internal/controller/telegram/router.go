package telegram

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gofiber/fiber/v2"
	"github.com/hong195/wheater-bot/internal/controller/telegram/request"
	"github.com/hong195/wheater-bot/internal/usecase/weather"
	"github.com/hong195/wheater-bot/pkg/logger"
)

func registerRoutes(app *fiber.App, bot *tgbotapi.BotAPI, w *weather.UseCase, l logger.Interface) {
	app.Post("/telegram/webhook", func(c *fiber.Ctx) error {
		ctx := c.Context()

		update := &tgbotapi.Update{}

		if err := c.BodyParser(update); err != nil {
			l.Error("Failed to parse update", "error", err)
			return c.Status(400).SendString("Bad Request")
		}

		if update.Message == nil {
			l.Info("Received update without message")
			return c.SendStatus(fiber.StatusOK)
		}

		handler := NewHandler(bot, update.Message.Chat.ID, w)
		err := handleCommands(ctx, update, handler)

		if err != nil {
			l.Error("Failed to handle command", "error", err)
		}

		return c.SendStatus(200)
	})
}

func handleCommands(ctx context.Context, update *tgbotapi.Update, h *Handler) error {
	switch update.Message.Command() {
	case StartCommand:
		return h.StartCommand()
	case GetWeatherDetailsCommand:

		coords, err := request.ParseCoordinates(update.Message.Text)

		if err != nil {
			return fmt.Errorf("error parsing coordinates: %w", err)
		}

		err = h.GetWeatherCommand(ctx, coords.Lat, coords.Lon)

		if err != nil {
			if err := h.SendATextMsg(err.Error()); err != nil {
				return fmt.Errorf("send message failed: %w", err)
			}
		}

		return fmt.Errorf("error getting weather: %w", err)
	default:
		return h.UnknownCommand()
	}
}
