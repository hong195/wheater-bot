package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type (
	// Config -.
	Config struct {
		App         App
		HTTP        HTTP
		Log         Log
		Metrics     Metrics
		Swagger     Swagger
		OpenWeather OpenWeather
		Telegram    Telegram
	}

	// App -.
	App struct {
		Name    string `env:"APP_NAME,required"`
		Version string `env:"APP_VERSION,required"`
	}

	// HTTP -.
	HTTP struct {
		Port           string `env:"HTTP_PORT,required"`
		UsePreforkMode bool   `env:"HTTP_USE_PREFORK_MODE" envDefault:"false"`
	}

	// Log -.
	Log struct {
		Level string `env:"LOG_LEVEL,required"`
	}

	// Metrics -.
	Metrics struct {
		Enabled bool `env:"METRICS_ENABLED" envDefault:"true"`
	}

	// Swagger -.
	Swagger struct {
		Enabled bool `env:"SWAGGER_ENABLED" envDefault:"true"`
	}

	OpenWeather struct {
		ApiUrl string `env:"OPEN_WEATHER_API_URL,required"`
		ApiKey string `env:"OPEN_WEATHER_API_KEY,required"`
	}

	Telegram struct {
		Token      string `env:"TELEGRAM_TOKEN,required"`
		WebhookURL string `env:"TELEGRAM_WEBHOOK_URL,required"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {

	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading config: %w", err)
	}

	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
