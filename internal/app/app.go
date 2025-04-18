package app

import (
	"log/slog"
	"mini-app-notifications/internal/bot"
	"mini-app-notifications/internal/config"
	"mini-app-notifications/internal/internal/service/event"
	"mini-app-notifications/internal/kafka"
	sl "mini-app-notifications/internal/logger"
	repository "mini-app-notifications/internal/repository/user"
	"mini-app-notifications/internal/storage/postgresql"
	"net/http"
	"os"
)

func Run(config *config.Config, logger *slog.Logger) {
	logger.Info("app running", slog.Attr{Key: "env", Value: slog.StringValue(config.Env)})

	// storage
	storage, err := postgresql.NewPostgreSQL(config.DB)
	if err != nil {
		logger.Error("failed to init postgresql storage", sl.Err(err))
		os.Exit(1)
	}

	// repo
	userRepo := repository.NewUserRepository(storage.DB, logger)

	// tgbot
	tgbot, err := bot.NewTelegramBot(config.TGBOT.BotToken)
	if err != nil {
		logger.Error("failed to create bot instanse", sl.Err(err))
		os.Exit(1)
	}

	// services
	eventSrv := event.NewEventService(logger, tgbot, userRepo)

	kafka.StartConsumer(config.Kafka, logger, eventSrv)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	http.ListenAndServe(config.Host+":"+config.Port, nil)
}
