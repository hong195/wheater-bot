// Package app configures and runs application.
package app

import (
	"fmt"
	"github.com/hong195/wheater-bot/internal/repo/webapi"
	"os"
	"os/signal"
	"syscall"

	"github.com/hong195/wheater-bot/config"
	"github.com/hong195/wheater-bot/internal/controller/http"
	"github.com/hong195/wheater-bot/internal/usecase/weather"
	"github.com/hong195/wheater-bot/pkg/httpserver"
	"github.com/hong195/wheater-bot/pkg/logger"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	weatherRepo := webapi.NewWeatherWebApi(cfg.OpenWeather.ApiUrl, cfg.OpenWeather.ApiKey)
	cityDetailRepo := webapi.NewCityDetailsRepository()
	// Use-Case
	weatherUseCase := weather.New(cityDetailRepo, weatherRepo)

	// HTTP Server
	httpServer := httpserver.New(httpserver.Port(cfg.HTTP.Port), httpserver.Prefork(cfg.HTTP.UsePreforkMode))
	http.NewRouter(httpServer.App, cfg, weatherUseCase, l)

	// Start servers
	httpServer.Start()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: %s", s.String())
	case err := <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err := httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
