// Package app configures and runs application.
package app

import (
	"fmt"
	"github.com/hong195/wheater-bot/internal/repo/webapi"
	stdhttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hong195/wheater-bot/config"
	"github.com/hong195/wheater-bot/internal/controller/http"
	"github.com/hong195/wheater-bot/internal/usecase/weather"
	"github.com/hong195/wheater-bot/pkg/httpserver"
	"github.com/hong195/wheater-bot/pkg/logger"
)

const reverseGeocodingUrl = "https://nominatim.openstreetmap.org/reverse"

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	lgr := logger.New(cfg.Log.Level)

	httpClient := &stdhttp.Client{
		Timeout: 3 * time.Second,
	}

	weatherRepo := webapi.NewWeatherWebApi(
		httpClient,
		cfg.OpenWeather.ApiKey,
		cfg.OpenWeather.ApiUrl,
	)

	cityDetailRepo := webapi.NewCityDetailsRepository(httpClient, reverseGeocodingUrl)
	// Use-Case
	weatherUseCase := weather.New(cityDetailRepo, weatherRepo)

	// HTTP Server
	httpServer := httpserver.New(httpserver.Port(cfg.HTTP.Port), httpserver.Prefork(cfg.HTTP.UsePreforkMode))
	http.NewRouter(httpServer.App, cfg, weatherUseCase, lgr)

	// Start servers
	httpServer.Start()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		lgr.Info("app - Run - signal: %s", s.String())
	case err := <-httpServer.Notify():
		lgr.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))

	}

	// Shutdown
	err := httpServer.Shutdown()
	if err != nil {
		lgr.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
