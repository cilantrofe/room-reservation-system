package config

// Настройки конфигурации

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Kafka struct {
		Broker string
		Topics []string
	}
	Telegram struct {
		Token string
	}
}

func LoadConfig() *Config {
	// Автоматически не подгружает, поэтому руками
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	cfg := &Config{}

	// Чтение переменных окружения
	cfg.Kafka.Broker = os.Getenv("KAFKA_BROKER")
	cfg.Kafka.Topics = []string{os.Getenv("KAFKA_TOPIC_CLIENT"), os.Getenv("KAFKA_TOPIC_HOTEL")}
	cfg.Telegram.Token = os.Getenv("TELEGRAM_TOKEN")

	return cfg
}
