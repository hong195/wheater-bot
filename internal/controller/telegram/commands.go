package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

const (
	StartCommand             = "start"
	GetWeatherDetailsCommand = "get_weather_details"
)

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
