package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Quizert/room-reservation-system/BookingSvc/internal/clients/grpc/hotelsvc"
	"github.com/Quizert/room-reservation-system/BookingSvc/internal/models"
	"github.com/Quizert/room-reservation-system/HotelSvc/api/grpc/hotelpb"
	"log"
	"time"
)

type BookingService struct {
	storage            Storage
	messageProducer    MessageProducer
	hotelSvcGrpcClient *grpc.HotelSvcClient
}

func NewBookingService(db Storage, producer MessageProducer, client *grpc.HotelSvcClient) *BookingService {
	return &BookingService{db, producer, client}
}

func (b *BookingService) GetBookingsByUserID(ctx context.Context, userID int) ([]*models.Booking, error) {
	return b.storage.GetBookingsByUserID(ctx, userID)
}

func (b *BookingService) GetBookingsByHotelID(ctx context.Context, id int) (*models.Booking, error) {
	return b.storage.GetBookingsByHotelID(ctx, id)
}

func (b *BookingService) CreateBooking(ctx context.Context, bookingRequest *models.BookingRequest) error {
	// Тут МБ валидация
	booking := bookingRequest.ToBooking()
	err := b.storage.CreateBooking(ctx, booking)

	if err != nil {
		return fmt.Errorf("error in CreateBooking: %w", err)
	}

	bookingMessage := bookingRequest.ToBookingMessage()
	kafkaMessageJSON, err := json.Marshal(bookingMessage)
	if err != nil {
		return fmt.Errorf("error in Marshal json: %w", err)
	}
	err = b.messageProducer.SendMessage(ctx, kafkaMessageJSON)
	if err != nil {
		log.Printf("Failed to send Kafka message: %v", err)
	}

	return nil
}

func (b *BookingService) UpdateBooking(ctx context.Context, booking *models.Booking) error {
	return b.storage.UpdateBooking(ctx, booking)
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
