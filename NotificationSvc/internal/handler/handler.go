package handler

// Обработка сообщений и передача их в сервис уведомлений

import (
	"NotificationSvc/internal/service"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
)

type BookingEvent struct {
	BookingID       int    `json:"booking_id"`
	HotelName       string `json:"hotel_name"`
	RoomDescription string `json:"room_description"`
	RoomNumber      int    `json:"room_number"`
	StartDate       string `json:"start_date"`
	EndDate         string `json:"end_date"`
	UserName        string `json:"user_name"`
	ChatId          string `json:"chat_id"`
}

type NotificationHandler struct {
	notificationService *service.NotificationService
}

func NewNotificationHandler(service *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: service,
	}
}

// Разбор сообщения для отеля
func (h *NotificationHandler) handleHotelBookingEvent(ctx context.Context, message string) {
	var event BookingEvent
	err := json.Unmarshal([]byte(message), &event)
	if err != nil {
		log.Printf("Failed to unmarshal hotel booking event: %v", err)
		return
	}

	chatId, _ := strconv.ParseInt(event.ChatId, 10, 64)

	notificationMessage := fmt.Sprintf(
		"Новое бронирование отеля:\n"+
			"Название отеля: %s\n"+
			"Описание номера: %s\n"+
			"Номер комнаты: %d\n"+
			"Дата заезда: %s\n"+
			"Дата выезда: %s\n"+
			"Имя гостя: %s\n",
		event.HotelName, event.RoomDescription, event.RoomNumber, event.StartDate, event.EndDate, event.UserName,
	)

	h.notificationService.SendNotification(notificationMessage, chatId)
}

// Разбор сообщения для клиента
func (h *NotificationHandler) handleClientBookingEvent(ctx context.Context, message string) {
	var event BookingEvent
	err := json.Unmarshal([]byte(message), &event)
	if err != nil {
		log.Printf("Failed to unmarshal client booking event: %v", err)
		return
	}

	chatId, _ := strconv.ParseInt(event.ChatId, 10, 64)

	notificationMessage := fmt.Sprintf(
		"%s, у Вас новое бронирование отеля.\nОзнакомьтесь с информацией ниже:\n"+
			"Название отеля: %s\n"+
			"Описание номера: %s\n"+
			"Номер комнаты: %d\n"+
			"Дата заезда: %s\n"+
			"Дата выезда: %s\n",
		event.UserName, event.HotelName, event.RoomDescription, event.RoomNumber, event.StartDate, event.EndDate,
	)

	h.notificationService.SendNotification(notificationMessage, chatId)
}

// Формирует уведомление о бронировании и отправляет его через NotificationService
func (h *NotificationHandler) HandleBookingEvent(ctx context.Context, message string, topic string) {
	switch topic {
	case "booking-topic-hotel":
		h.handleHotelBookingEvent(ctx, message)
	case "booking-topic-client":
		h.handleClientBookingEvent(ctx, message)
	default:
		log.Printf("Unknown topic: %s", topic)
	}
}
