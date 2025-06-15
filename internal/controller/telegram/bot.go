package telegram

import (
	"context"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gofiber/fiber/v2"
	"github.com/hong195/wheater-bot/config"
	"github.com/hong195/wheater-bot/internal/controller/http/middleware"
	"github.com/hong195/wheater-bot/internal/usecase/weather"
	"github.com/hong195/wheater-bot/pkg/logger"
	"strconv"
	"strings"
)

const (
	StartCommand             = "start"
	GetWeatherDetailsCommand = "get_weather_details"
)

func NewBot(app *fiber.App, cfg *config.Config, w *weather.UseCase, l logger.Interface) {

	app.Use(middleware.Recovery(l))

	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.Token)
	if err != nil {
		l.Fatal("Failed to create bot API", "error", err)
	}

	if err := registerCommands(bot); err != nil {
		l.Fatal("Error happened while registering command", err)
	}

	if err := registerWebhook(bot, cfg.Telegram.WebhookURL); err != nil {
		l.Fatal("Error happened while registering command", err)
	}

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

		lat, lon, err := parseCoordinates(update.Message.Text)

		if err != nil {
			return fmt.Errorf("error parsing coordinates: %w", err)
		}

		// https://nominatim.openstreetmap.org/reverse?lat=41.31053&lon=69.24056&lang=ru&format=jsonv2
		return h.GetWeatherCommand(ctx, lat, lon)
	default:
		return h.UnknownCommand()
	}
}

func parseCoordinates(text string) (float64, float64, error) {
	// Удаляем команду
	parts := strings.SplitN(text, " ", 2)
	if len(parts) != 2 {
		return 0, 0, errors.New("incorrect coords format")
	}

	coords := strings.Split(parts[1], ",")
	if len(coords) != 2 {
		return 0, 0, errors.New("coords should be separated by comma")
	}

	lat, err1 := strconv.ParseFloat(strings.TrimSpace(coords[0]), 64)
	lon, err2 := strconv.ParseFloat(strings.TrimSpace(coords[1]), 64)

	if err1 != nil || err2 != nil {
		return 0, 0, errors.New("error parsing coordinates")
	}

	return lat, lon, nil
}

func registerWebhook(bot *tgbotapi.BotAPI, webhookUrl string) error {
	wh, err := tgbotapi.NewWebhook(webhookUrl)
	if err != nil {
		return fmt.Errorf("create webhook: %w", err)
	}
	if _, err := bot.Request(wh); err != nil {
		return fmt.Errorf("set webhook: %w", err)
	}
	return nil
}

func registerCommands(bot *tgbotapi.BotAPI) error {
	commands := []tgbotapi.BotCommand{
		{
			Command:     StartCommand,
			Description: "Запустить бота",
		},
		{
			Command:     GetWeatherDetailsCommand,
			Description: "Получить погоду по координатам",
		},
	}

	cfg := tgbotapi.NewSetMyCommands(commands...)
	_, err := bot.Request(cfg)

	return err
}
