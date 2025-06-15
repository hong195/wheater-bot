package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gofiber/fiber/v2"
	"github.com/hong195/wheater-bot/config"
	"github.com/hong195/wheater-bot/internal/controller/http/middleware"
	"github.com/hong195/wheater-bot/internal/usecase/weather"
	"github.com/hong195/wheater-bot/pkg/logger"
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
