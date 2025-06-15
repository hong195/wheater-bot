package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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
