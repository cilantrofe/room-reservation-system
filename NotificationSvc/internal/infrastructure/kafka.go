package infrastructure

// Работа с Kafka

import (
	"NotificationSvc/internal/handler"
	"context"
	"github.com/segmentio/kafka-go"
	"log"
)

type KafkaConsumer struct {
	readers             []*kafka.Reader
	notificationHandler *handler.NotificationHandler
}

// Инициализация KafkaConsumer с конфигурацией и хэндлером
func NewKafkaConsumer(broker string, topics []string, notificationHandler *handler.NotificationHandler) (*KafkaConsumer, error) {
	var readers []*kafka.Reader
	for _, topic := range topics {
		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{broker},
			Topic:   topic,
			GroupID: "notification-group",
		})
		readers = append(readers, r)
	}

	return &KafkaConsumer{
		readers:             readers,
		notificationHandler: notificationHandler,
	}, nil
}

// Запуск процесса чтения сообщений из Kafka
func (kc *KafkaConsumer) StartConsuming(ctx context.Context) {
	for _, reader := range kc.readers {
		go func(r *kafka.Reader) {
			for {
				m, err := r.ReadMessage(ctx)
				if err != nil {
					log.Printf("Error reading message from Kafka: %v", err)
					continue
				}
				// Передаём сообщение в обработчик уведомлений
				kc.notificationHandler.HandleBookingEvent(ctx, string(m.Value), m.Topic)
			}
		}(reader)
	}
}
