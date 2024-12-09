package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Quizert/room-reservation-system/BookingSvc/internal/clients/grpc/hotelsvc"
	"github.com/Quizert/room-reservation-system/BookingSvc/internal/clients/http/paymentsvc"
	"github.com/Quizert/room-reservation-system/BookingSvc/internal/models"
	"github.com/Quizert/room-reservation-system/HotelSvc/api/grpc/hotelpb"
	"log"
	"time"
)

type BookingService struct {
	storage             Storage
	messageProducer     MessageProducer
	hotelSvcGrpcClient  *grpc.HotelSvcClient
	paymentSystemClient *paymentsvc.Client
}

func NewBookingService(db Storage, producer MessageProducer, hotelClient *grpc.HotelSvcClient, paymentClient *paymentsvc.Client) *BookingService {
	return &BookingService{db, producer, hotelClient, paymentClient}
}

func (b *BookingService) CreateBooking(ctx context.Context, bookingRequest *models.BookingRequest) error {
	// Тут МБ валидация
	booking := bookingRequest.ToBooking()

	bookingID, err := b.storage.CreateBooking(ctx, booking)
	if err != nil {
		return fmt.Errorf("error in CreateBooking: %w", err)
	}

	bookingMessage := bookingRequest.ToBookingMessage(bookingID)
	paymentRequest := paymentsvc.ToPaymentRequest(bookingMessage, bookingRequest.CardNumber, bookingRequest.Amount)

	err = b.paymentSystemClient.CreatePaymentRequest(ctx, paymentRequest)
	if err != nil {
		return fmt.Errorf("error in payment request: %w", err)
	}
	return nil
}

func (b *BookingService) GetBookingsByUserID(ctx context.Context, userID int) ([]*models.Booking, error) {
	return b.storage.GetBookingsByUserID(ctx, userID)
}

func (b *BookingService) GetBookingsByHotelID(ctx context.Context, id int) (*models.Booking, error) {
	return b.storage.GetBookingsByHotelID(ctx, id)
}

func (b *BookingService) UpdateBookingStatus(ctx context.Context, status string, bookingMessage *models.BookingMessage) error {
	err := b.storage.UpdateBookingStatus(ctx, status, bookingMessage.BookingID)
	if err != nil {
		return fmt.Errorf("error in UpdateBookingStatus: %w", err)
	}

	switch status {
	case "success":
		kafkaUserMessage, err := json.Marshal(bookingMessage)
		if err != nil {
			return fmt.Errorf("error in Marshal KafkaUserMessage: %w", err)
		}
		err = b.messageProducer.SendMessage(ctx, kafkaUserMessage)
		if err != nil {
			return fmt.Errorf("error in SendMessage: %w", err)
		}

		hotelierMessage := bookingMessage.ToHotelierMessage("hotelier name", "123123")

		kafkaHotelierMessage, err := json.Marshal(hotelierMessage)
		if err != nil {
			return fmt.Errorf("error in Marshal KafkaHotelierMessage: %w", err)
		}
		err = b.messageProducer.SendMessage(ctx, kafkaHotelierMessage)
		if err != nil {
			return fmt.Errorf("error in SendMessage: %w", err)
		}

	case "fail":
		log.Println("failed")
	}
	return nil
}

func (b *BookingService) DeleteBooking(ctx context.Context, id int) error {
	return b.storage.DeleteBooking(ctx, id)
}

func (b *BookingService) GetAvailableRooms(ctx context.Context, hotelID int, startDate, endDate time.Time) ([]*hotelpb.Room, error) {
	request := hotelpb.GetRoomsRequest{HotelId: int32(hotelID)}
	allRooms, err := b.hotelSvcGrpcClient.Api.GetRoomsByHotelId(ctx, &request)
	if err != nil {
		return nil, fmt.Errorf("error in gRPC request GetRoomsByHotelID: %v", err)
	}
	log.Println("ALL ROOMS:", allRooms)
	unavailableRoomsID, err := b.storage.GetUnavailableRoomsByHotelId(ctx, hotelID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error in db GetUnavailableRoomsByHotelId: %v", err)
	}
	log.Println("UNAVAILABLE ROOMS:", unavailableRoomsID)
	availableRooms := make([]*hotelpb.Room, 0)
	for _, room := range allRooms.Rooms {
		if _, ok := unavailableRoomsID[int(room.Id)]; !ok {
			availableRooms = append(availableRooms, room)
		}
	}
	return availableRooms, nil
}
