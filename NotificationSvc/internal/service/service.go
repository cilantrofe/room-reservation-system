package service

// Отправка уведомлений отелям и клиентам через Telegram, используя ChatID

import "NotificationSvc/internal/delivery"

type NotificationService struct {
	notifier *delivery.TelegramNotifier
}

func NewNotificationService(notifier *delivery.TelegramNotifier) *NotificationService {
	return &NotificationService{notifier: notifier}
}

// Наконец-то добираемся до отправки сообщения через Telegram
func (s *NotificationService) SendNotification(message string) error {
	return s.notifier.SendNotification(message)
}
