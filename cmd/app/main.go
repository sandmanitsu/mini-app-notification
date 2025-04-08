package main

import (
	"log"
	"mini-app-notifications/internal/app"
	"mini-app-notifications/internal/config"
	"mini-app-notifications/internal/logger"
)

func main() {
	log.Println("config initializing...")
	config := config.MustLoad()

	log.Println("logger initializing...")
	logger := logger.NewLogger(config.Env)

	app.Run(config, logger)
}
