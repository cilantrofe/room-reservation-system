package infrastructure

// Работа с Kafka

import (
	"NotificationSvc/internal/handler"
	"context"
	"github.com/segmentio/kafka-go"
	"log"
)

type KafkaConsumer struct {
	reader              *kafka.Reader
	notificationHandler *handler.NotificationHandler
}

// Инициализация KafkaConsumer с конфигурацией и хэндлером
func NewKafkaConsumer(broker, topic string, notificationHandler *handler.NotificationHandler) (*KafkaConsumer, error) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   topic,
		GroupID: "notification-group", // если не устанавливать, то все консьюмеры будут читать одно и то же
	})

	return &KafkaConsumer{
		reader:              r,
		notificationHandler: notificationHandler,
	}, nil
}

// Запуск процесса чтения сообщений из Kafka
func (kc *KafkaConsumer) StartConsuming(ctx context.Context) {
	for {
		m, err := kc.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Error reading message from Kafka: %v", err)
			continue
		}

		// Передаём сообщение в обработчик уведомлений
		kc.notificationHandler.HandleBookingEvent(ctx, string(m.Value))
	}
}
