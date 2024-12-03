package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) *Producer {
	return &Producer{
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers:      brokers,
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			BatchTimeout: 10 * time.Microsecond,
		}),
	}
}

func (p *Producer) SendMessage(ctx context.Context, value []byte) error {
	message := kafka.Message{
		Value: value,
	}
	err := p.writer.WriteMessages(ctx, message)
	if err != nil {
		log.Printf("Failed to send Kafka message: %v", err)
		return fmt.Errorf("failed to send Kafka message: %w", err)
	}
	log.Printf("Message sent to Kafka: %s", value)
	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
