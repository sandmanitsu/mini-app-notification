package kafka

import (
	"context"
	"fmt"
	"log/slog"
	"mini-app-notifications/internal/config"
	sl "mini-app-notifications/internal/logger"

	"github.com/segmentio/kafka-go"
)

func StartConsumer(cfg config.Kafka, logger *slog.Logger) {
	const op = "kafka.StartConsumer"

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{cfg.KafkaBroker},
		Topic:     cfg.KafkaTopic,
		Partition: 0,
		MaxBytes:  10e6,
	})
	ctx := context.Background()

	logger.Info("starting read kafka")
	go func() {
		for {
			m, err := r.ReadMessage(ctx)
			if err != nil {
				logger.Error(fmt.Sprintf("%s, read message failed", op), sl.Err(err))
			}

			fmt.Println(string(m.Value))
		}
	}()
}
