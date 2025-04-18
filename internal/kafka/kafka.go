package kafka

import (
	"context"
	"fmt"
	"log/slog"
	"mini-app-notifications/internal/config"
	"mini-app-notifications/internal/domain"
	sl "mini-app-notifications/internal/logger"
	"time"

	"github.com/segmentio/kafka-go"
)

type EventService interface {
	Process(event domain.Event) error
}

func StartConsumer(cfg config.Kafka, logger *slog.Logger, eventSrv EventService) {
	const op = "kafka.StartConsumer"

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{cfg.KafkaBroker},
		Topic:   cfg.KafkaTopic,
		GroupID: "notification-service",
		// Partition:   0,
		MaxBytes:    10e6,
		Logger:      kafka.LoggerFunc(logf),
		ErrorLogger: kafka.LoggerFunc(logf),
		Dialer: &kafka.Dialer{
			Timeout:   10 * time.Second,
			DualStack: true,
			ClientID:  "my-client-v1",
		},
	})
	ctx := context.Background()

	logger.Info("starting read kafka...")
	go func() {
		for {
			m, err := r.FetchMessage(ctx)
			if err != nil {
				logger.Error(fmt.Sprintf("%s : failed fetch message", op), sl.Err(err))
				continue
			}

			if err = eventSrv.Process(domain.Event{EventType: string(m.Key), Value: m.Value}); err != nil {
				logger.Error(fmt.Sprintf("%s : failed process notification message", op), sl.Err(err))
				continue
			}

			if err := r.CommitMessages(ctx, m); err != nil {
				logger.Error(fmt.Sprintf("%s : failed commit message", op), sl.Err(err))
			}
		}
	}()
}

func logf(msg string, a ...interface{}) {
	fmt.Printf(msg, a...)
	fmt.Println()
}
