package handler

// Обработка сообщений и передача их в сервис уведомлений

import (
	"NotificationSvc/internal/service"
	"context"
)

type NotificationHandler struct {
	notificationService *service.NotificationService
}

func NewNotificationHandler(service *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: service,
	}
}

// Формирует уведомление о бронировании и отправляет его через NotificationService
func (h *NotificationHandler) HandleBookingEvent(ctx context.Context, message string) {
	// TODO: Разобрать полученный json от кафки producer на инфу о клиенте/отеле, времени бронирования и самом бронировании
	h.notificationService.SendNotification("Новое бронирование: " + message)
}
