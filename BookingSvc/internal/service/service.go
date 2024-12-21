package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Quizert/room-reservation-system/BookingSvc/internal/models"
	"github.com/Quizert/room-reservation-system/HotelSvc/api/grpc/hotelpb"
	"go.uber.org/zap"
	"log"
	"time"
)

var (
	ErrForbiddenAccess = errors.New("forbidden access")
)

type BookingServiceImpl struct {
	storage             Storage
	messageProducer     MessageProducer
	hotelSvcClient      HotelClient
	paymentSystemClient PaymentSystemClient
	log                 *zap.Logger
}

func NewBookingServiceImpl(db Storage, producer MessageProducer, hotelClient HotelClient, paymentClient PaymentSystemClient, logger *zap.Logger) *BookingServiceImpl {
	return &BookingServiceImpl{db, producer, hotelClient, paymentClient, logger}
}

func (b *BookingServiceImpl) CreateBooking(ctx context.Context, bookingRequest *models.BookingRequest, user *models.User) error {
	b.log.With(
		zap.String("Layer", "service: RegisterUser"),
		zap.Int("room id", bookingRequest.RoomID),
		zap.Int("is_hotelier", bookingRequest.HotelID),
		zap.Time("start date", bookingRequest.StartDate),
		zap.Time("end date", bookingRequest.EndDate),
		zap.String("hotel name", bookingRequest.HotelName),
		zap.String("RoomDescription", bookingRequest.RoomDescription),
		zap.String("card number", bookingRequest.CardNumber),
		zap.Int("user id", user.UserID),
		zap.String("username", user.Username),
		zap.String("chat id", user.ChatID),
		zap.Int("amount", bookingRequest.Amount)).Info("Received request to create booking")

	fmt.Println(user.UserID, user.Username, user.ChatID)

	booking := bookingRequest.ToBookingInfo(user.UserID)

	bookingID, err := b.storage.CreateBooking(ctx, booking)
	if err != nil {
		return fmt.Errorf("error in CreateBooking: %w", err)
	}

	//TODO: Расчитать Amount
	bookingMessage := bookingRequest.ToBookingMessage(bookingID, user.Username, user.ChatID)
	paymentRequest := models.ToPaymentRequest(bookingMessage, bookingRequest.CardNumber, bookingRequest.Amount)

	err = b.paymentSystemClient.CreatePaymentRequest(ctx, paymentRequest)
	if err != nil {
		return fmt.Errorf("error in payment request: %w", err)
	}
	return nil
}

func (b *BookingServiceImpl) GetBookingsByUserID(ctx context.Context, userID int) ([]*models.BookingInfo, error) {
	return b.storage.GetBookingsByUserID(ctx, userID)
}

func (b *BookingServiceImpl) GetBookingsByHotelID(ctx context.Context, hotelID, userID int) (*models.BookingInfo, error) {
	//Проверим, что отель действительно принадлежит userID
	b.hotelSvcClient.GetRoomsByHotelId()
	return b.storage.GetBookingsByHotelID(ctx, id)
}

func (b *BookingServiceImpl) UpdateBookingStatus(ctx context.Context, status string, bookingMessage *models.BookingMessage) error {
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
		err = b.messageProducer.SendUserMessage(ctx, kafkaUserMessage)
		if err != nil {
			return fmt.Errorf("error in SendMessage: %w", err)
		}
		log.Println("message sent", bookingMessage)
		//TODO: gRPC к AuthSvc для HotelierChatID

		hotelierMessage := bookingMessage.ToHotelierMessage("hotelier name", "123123")

		kafkaHotelierMessage, err := json.Marshal(hotelierMessage)
		if err != nil {
			return fmt.Errorf("error in Marshal KafkaHotelierMessage: %w", err)
		}
		err = b.messageProducer.SendHotelierMessage(ctx, kafkaHotelierMessage)
		if err != nil {
			return fmt.Errorf("error in SendMessage: %w", err)
		}

	case "fail":
		log.Println("failed")
	}
	return nil
}

func (b *BookingServiceImpl) DeleteBooking(ctx context.Context, id int) error {
	return b.storage.DeleteBooking(ctx, id)
}

func (b *BookingServiceImpl) GetAvailableRooms(ctx context.Context, hotelID int, startDate, endDate time.Time) ([]*hotelpb.Room, error) {
	request := hotelpb.GetRoomsRequest{HotelId: int32(hotelID)}
	allRooms, err := b.hotelSvcClient.GetRoomsByHotelId(ctx, &request)
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
