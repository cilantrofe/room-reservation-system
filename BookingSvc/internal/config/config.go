package config

import (
	"os"
)

type Config struct {
	DBHost           string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBName           string
	GRPCHost         string
	GRPCPort         string
	HTTPPort         string
	KafkaBroker      string
	KafkaTopicClient string
	KafkaTopicHotel  string
	PaymentSvcURL    string
}

func LoadConfig() (*Config, error) {
	return &Config{
		DBHost:           os.Getenv("BOOKING_DB_HOST"),
		DBPort:           os.Getenv("BOOKING_DB_PORT"),
		DBUser:           os.Getenv("BOOKING_DB_USER"),
		DBPassword:       os.Getenv("BOOKING_DB_PASSWORD"),
		DBName:           os.Getenv("BOOKING_DB_NAME"),
		GRPCHost:         os.Getenv("BOOKING_GRPC_HOST"),
		GRPCPort:         os.Getenv("BOOKING_GRPC_PORT"),
		HTTPPort:         os.Getenv("BOOKING_HTTP_PORT"),
		KafkaBroker:      os.Getenv("KAFKA_BROKER"),
		KafkaTopicClient: os.Getenv("KAFKA_TOPIC_CLIENT"),
		KafkaTopicHotel:  os.Getenv("KAFKA_TOPIC_HOTEL"),
		PaymentSvcURL:    os.Getenv("PAYMENT_SERVICE_URL"),
	}, nil
}
