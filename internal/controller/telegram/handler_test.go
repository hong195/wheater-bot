package telegram

import (
	"context"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hong195/wheater-bot/internal/repo/webapi"
	"github.com/hong195/wheater-bot/internal/usecase/weather"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

const (
	lat = 41.311081
	lon = 69.240562
)

type mockBotAPI struct {
	mock.Mock
}

func (m *mockBotAPI) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	args := m.Called(c)
	return args.Get(0).(tgbotapi.Message), args.Error(1)
}

func getUseCase() *weather.UseCase {

	cityDetails := webapi.NewInMemoryCityRepo()
	weatherRepo := webapi.NewInMemoryWeatherRepo()

	return weather.NewUseCase(cityDetails, weatherRepo)
}

func TestGetWeatherCommand(t *testing.T) {
	mockBot := new(mockBotAPI)
	chatID := int64(12345)

	useCase := getUseCase()
	handler := NewHandler(mockBot, chatID, useCase)

	expectedMsg := tgbotapi.NewMessage(chatID, mock.Anything)
	expectedMsg.ParseMode = tgbotapi.ModeMarkdown

	mockBot.
		On("Send", mock.AnythingOfType("tgbotapi.MessageConfig")).
		Return(tgbotapi.Message{MessageID: 1}, nil)

	err := handler.GetWeatherCommand(context.Background(), lat, lon)

	assert.NoError(t, err)
	mockBot.AssertExpectations(t)
}

func TestGetWeatherCommand_Error(t *testing.T) {
	// Arrange
	mockBot := new(mockBotAPI)
	chatID := int64(12345)

	useCase := getUseCase()
	handler := NewHandler(mockBot, chatID, useCase)

	mockBot.
		On("Send", mock.AnythingOfType("tgbotapi.MessageConfig")).
		Return(tgbotapi.Message{MessageID: 1}, errors.New("some error"))

	// Act
	err := handler.GetWeatherCommand(context.Background(), lat, lon)

	// Assert
	assert.Error(t, err)
	mockBot.AssertExpectations(t)
}
