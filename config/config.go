package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type (
	// Config -.
	Config struct {
		App         App
		HTTP        HTTP
		Log         Log
		Metrics     Metrics
		Swagger     Swagger
		NGROK       NGROK
		OpenWeather OpenWeather
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
		Enabled bool `env:"SWAGGER_ENABLED" envDefault:"false"`
	}

	NGROK struct {
		AuthToken string `env:"NGROK_AUTHTOKEN" envDefault:"false"`
	}

	OpenWeather struct {
		ApiUrl string `env:"OPEN_WEATHER_API_URL" envDefault:"false"`
		ApiKey string `env:"OPEN_WEATHER_API_KEY" envDefault:"false"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
