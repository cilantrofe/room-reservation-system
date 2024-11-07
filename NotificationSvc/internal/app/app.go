package app

// Инцилизация и запуск компонентов приложения

import (
	"NotificationSvc/internal/config"
	"NotificationSvc/internal/delivery"
	"NotificationSvc/internal/handler"
	"NotificationSvc/internal/infrastructure"
	"NotificationSvc/internal/service"
	"context"
)

type App struct {
	kafkaConsumer *infrastructure.KafkaConsumer
}

func NewApp() *App {
	return &App{}
}

func (a *App) Init(ctx context.Context) error {
	cfg := config.LoadConfig()

	// Создание TelegramNotifier
	telegramNotifier, err := delivery.NewTelegramNotifier(cfg.Telegram.Token, cfg.Telegram.ChatID)
	if err != nil {
		return err
	}

	// Инициализация NotificationService и NotificationHandler
	notificationService := service.NewNotificationService(telegramNotifier)
	notificationHandler := handler.NewNotificationHandler(notificationService)

	// Инициализация KafkaConsumer с конфигурацией и хэндлером
	kafkaConsumer, err := infrastructure.NewKafkaConsumer(cfg.Kafka.Broker, cfg.Kafka.Topic, notificationHandler)
	if err != nil {
		return err
	}
	a.kafkaConsumer = kafkaConsumer

	return nil
}

func (a *App) Start(ctx context.Context) error {
	go a.kafkaConsumer.StartConsuming(ctx)
	return nil
}
