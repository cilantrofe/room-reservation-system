package config

// Настройки конфигурации

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type Config struct {
	Kafka struct {
		Broker string
		Topic  string
	}
	Telegram struct {
		Token  string
		ChatID int64
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
	cfg.Kafka.Topic = os.Getenv("KAFKA_TOPIC")
	cfg.Telegram.Token = os.Getenv("TELEGRAM_TOKEN")
	cfg.Telegram.ChatID, _ = strconv.ParseInt(os.Getenv("TELEGRAM_CHAT_ID"), 10, 64)

	// Для отладки
	log.Printf("Telegram Token: %s, Chat ID: %d", cfg.Telegram.Token, cfg.Telegram.ChatID)

	return cfg
}
