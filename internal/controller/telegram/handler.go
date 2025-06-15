package telegram

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hong195/wheater-bot/internal/usecase/weather"
	"time"
)

type Sender interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
}

type Handler struct {
	botApi  Sender
	chatID  int64
	useCase *weather.UseCase
}

func NewHandler(botApi Sender, chatID int64, useCase *weather.UseCase) *Handler {
	return &Handler{
		botApi:  botApi,
		chatID:  chatID,
		useCase: useCase,
	}
}

func (h *Handler) StartCommand() error {

	msg := tgbotapi.NewMessage(h.chatID, "Привет бот запущен! Отправь долготу и широту чтобы получить детали прогноза погоды")

	_, err := h.botApi.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) SendATextMsg(msg string) error {
	m := tgbotapi.NewMessage(h.chatID, msg)

	_, err := h.botApi.Send(m)

	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) UnknownCommand() error {
	msg := tgbotapi.NewMessage(h.chatID, "Неизвестная команда!")

	_, err := h.botApi.Send(msg)

	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) GetWeatherCommand(context context.Context, lat, lon float64) error {

	res, err := h.useCase.GetWeatherByCoordinates(context, lat, lon)

	if err != nil {
		return err
	}

	formattedStr, err := h.formatWeatherResponse(res)

	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(h.chatID, formattedStr)
	msg.ParseMode = tgbotapi.ModeMarkdown

	_, err = h.botApi.Send(msg)

	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) formatWeatherResponse(w *weather.Weather) (string, error) {

	if w == nil {
		return "", fmt.Errorf("weather data is nil")
	}

	msg := fmt.Sprintf(
		"🌤 *Погода сейчас*\n\n"+
			"📍 *Город*: %s\n"+
			"🗺 *Район*: %s\n"+
			"🏙 *Улица*: %s\n\n"+
			"🌅 *Рассвет*: %s\n"+
			"🌇 *Закат*: %s\n"+
			"🌡 *Температура*: %.1f°C\n"+
			"🤒 *Ощущается как*: %.1f°C\n"+
			"💧 *Влажность*: %d%%\n"+
			"💨 *Ветер*: %.1f м/с\n"+
			"🔎 *Описание*: %s",
		w.City,
		w.County,
		w.Neighborhood,
		time.Unix(int64(w.Sunrise), 0).Format("15:04"),
		time.Unix(int64(w.Sunset), 0).Format("15:04"),
		w.Temperature,
		w.FeelsLike,
		w.Humidity,
		w.WindSpeed,
		w.Description,
	)

	return msg, nil
}
