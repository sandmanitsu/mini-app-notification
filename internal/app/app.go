package app

import (
	"log/slog"
	"mini-app-notifications/internal/config"
	"mini-app-notifications/internal/kafka"
	"net/http"
)

func Run(config *config.Config, logger *slog.Logger) {
	logger.Info("app running", slog.Attr{Key: "env", Value: slog.StringValue(config.Env)})

	kafka.StartConsumer(config.Kafka, logger)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	http.ListenAndServe(config.Host+":"+config.Port, nil)
}
