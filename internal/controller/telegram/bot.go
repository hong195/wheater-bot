package telegram

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gofiber/fiber/v2"
	"github.com/hong195/wheater-bot/config"
	"github.com/hong195/wheater-bot/internal/controller/http/middleware"
	"github.com/hong195/wheater-bot/internal/usecase/weather"
	"github.com/hong195/wheater-bot/pkg/logger"
	"strconv"
	"strings"
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

	registerRoutes(app, bot, w, l)
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
