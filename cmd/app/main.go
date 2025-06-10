package main

import (
	"log"

	"github.com/hong195/wheater-bot/config"
	"github.com/hong195/wheater-bot/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app.Run(cfg)
}
