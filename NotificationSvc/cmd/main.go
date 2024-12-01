package main

// Инициализация и запуск микросервиса

import (
	"NotificationSvc/internal/app"
	"context"
	"log"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Инициализация приложения
	application := app.NewApp()
	if err := application.Init(ctx); err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	// Запуск приложения
	if err := application.Start(ctx); err != nil {
		log.Fatalf("Failed to start app: %v", err)
	}
}
