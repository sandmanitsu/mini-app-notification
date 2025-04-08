package config

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Host  string `env:"HOST" env-required:"true"`
	Port  string `env:"PORT" env-required:"true"`
	Env   string `env:"ENV" env-required:"true"`
	Kafka Kafka
}

type Kafka struct {
	KafkaBroker string `env:"KAFKA_BROKER" env-required:"true"`
	KafkaTopic  string `env:"KAFKA_TOPIC" env-required:"true"`
}

var (
	config *Config
	once   sync.Once
)

// Getting config variables from .env file
func MustLoad() *Config {
	if config == nil {
		once.Do(
			func() {
				configPath := filepath.Join(".env")

				if _, err := os.Stat(configPath); err != nil {
					log.Fatalf("Error opening config file: %s", err)
				}

				var newConfig Config
				err := cleanenv.ReadConfig(configPath, &newConfig)
				if err != nil {
					log.Fatalf("Error reading config file: %s", err)
				}

				config = &newConfig
			})
	}

	return config
}
