package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Quizert/room-reservation-system/AuthSvc/pkj/authpb"
	"github.com/Quizert/room-reservation-system/BookingSvc/internal/models"
	"github.com/Quizert/room-reservation-system/BookingSvc/internal/myerror"
	"github.com/Quizert/room-reservation-system/HotelSvc/api/grpc/hotelpb"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"time"
)

type BookingServiceImpl struct {
	storage             Storage
	messageProducer     MessageProducer
	hotelSvcClient      HotelClient
	authSvcClient       AuthSvcClient
	paymentSystemClient PaymentSystemClient
	log                 *zap.Logger
}

type AuthSvcClient interface {
	GetHotelierInformation(ctx context.Context, request *authpb.GetHotelierRequest) (*authpb.GetHotelierResponse, error)
}

func NewBookingServiceImpl(db Storage, producer MessageProducer, hotelClient HotelClient, authClient AuthSvcClient, paymentClient PaymentSystemClient, logger *zap.Logger) *BookingServiceImpl {
	return &BookingServiceImpl{storage: db, messageProducer: producer, hotelSvcClient: hotelClient, authSvcClient: authClient, paymentSystemClient: paymentClient, log: logger}
}

func (b *BookingServiceImpl) CreateBooking(ctx context.Context, bookingRequest *models.BookingRequest, user *models.User) error {
	bookingRequest.Amount = bookingRequest.CountOfPeople * bookingRequest.RoomBasePrice
	b.log.With(
		zap.String("Layer", "service: CreateBooking"),
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

	booking := bookingRequest.ToBookingInfo(user.UserID)

	bookingID, err := b.storage.CreateBooking(ctx, booking)
	if err != nil {
		if errors.Is(err, myerror.ErrBookingAlreadyExists) {
			b.log.Warn("in service Create Booking", zap.Error(err))
			return fmt.Errorf("in service Create Booking: %w", myerror.ErrBookingAlreadyExists)
		}
		b.log.Error("in service Create Booking", zap.Error(err))
		return fmt.Errorf("in service Create Booking: %w", err)
	}

	bookingMessage := bookingRequest.ToBookingMessage(bookingID, user.Username, user.ChatID)
	paymentRequest := models.ToPaymentRequest(bookingMessage, bookingRequest.CardNumber, bookingRequest.Amount)

	err = b.paymentSystemClient.CreatePaymentRequest(ctx, paymentRequest)
	if err != nil {
		return fmt.Errorf("error in payment request: %w", err)
	}

	b.log.Info("in service Create Booking end successfully")
	return nil
}

func (b *BookingServiceImpl) GetBookingsByUserID(ctx context.Context, userID int) ([]*models.BookingInfo, error) {
	b.log.With(
		zap.String("Layer", "service: GetBookingsByUserID"),
		zap.Int("user id", userID),
	).Info("Received request to get booking by user id")

	bookings, err := b.storage.GetBookingsByUserID(ctx, userID)
	if err != nil {
		b.log.Error("error in service GetBookingsByHotelID:", zap.Error(err))
		return nil, fmt.Errorf("error in service GetBookingsByHotelID: %w", err)
	}

	b.log.Info("in service get bookings by user id end successfully")
	return bookings, nil
}

func (b *BookingServiceImpl) GetBookingsByHotelID(ctx context.Context, hotelID, userID int) ([]*models.BookingInfo, error) {
	b.log.With(
		zap.String("Layer", "service: GetBookingsByHotelID"),
		zap.Int("user id", userID),
		zap.Int("hotel id", hotelID)).Info("Received request to get bookings by owner")

	//Проверим, что отель действительно принадлежит userID
	req := &hotelpb.GetOwnerIdRequest{Id: int32(hotelID)}
	response, err := b.hotelSvcClient.GetOwnerIdByHotelId(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			b.log.Warn("error in service gRPC GetBookingsByHotelID:", zap.Error(myerror.ErrHotelNotFound))
			return nil, fmt.Errorf("error in service GetOwnerIdByHotelId: %w", myerror.ErrHotelNotFound)
		}
		b.log.Error("error in service gRPC GetBookingsByHotelID:", zap.Error(err))
		return nil, fmt.Errorf("error in service GetOwnerIdByHotelId: %w", err)
	}
	if response.OwnerId != int32(userID) {
		b.log.Warn("hotelier tries to check forbidden bookings from other hotel!")
		return nil, fmt.Errorf("user not owned by hotel %w", myerror.ErrForbiddenAccess)
	}

	bookings, err := b.storage.GetBookingsByHotelID(ctx, hotelID)
	if err != nil {
		b.log.Error("error in service GetBookingsByHotelID:", zap.Error(err))
		return nil, fmt.Errorf("error in service GetBookingsByHotelID: %w", err)
	}

	b.log.Info("in service get bookings by hotel id end successfully")
	return bookings, nil
}

func (b *BookingServiceImpl) UpdateBookingStatus(ctx context.Context, BookingStatus string, bookingMessage *models.BookingMessage) error {
	b.log.With(
		zap.String("Layer", "service: UpdateBookingStatus"),
		zap.String("status", BookingStatus),
		zap.Any("message", bookingMessage),
	).Info("Received request to update booking status")

	err := b.storage.UpdateBookingStatus(ctx, BookingStatus, bookingMessage.BookingID)
	if err != nil {
		b.log.Error("error in service UpdateBookingStatus: %w", zap.Error(err))
		return fmt.Errorf("error in service UpdateBookingStatus: %w", err)
	}

	switch BookingStatus {
	case "success":
		kafkaUserMessage, err := json.Marshal(bookingMessage)
		if err != nil {
			b.log.Error("error in service UpdateBookingStatus", zap.Error(err))
			return fmt.Errorf("error in Marshal KafkaUserMessage: %w", err)
		}
		err = b.messageProducer.SendUserMessage(ctx, kafkaUserMessage)
		if err != nil {
			b.log.Error("error in service UpdateBookingStatus", zap.Error(err))
			return fmt.Errorf("error SendMessage: %w", err)
		}

		hotelReq := &hotelpb.GetOwnerIdRequest{Id: int32(bookingMessage.HotelID)}
		hotelResponse, err := b.hotelSvcClient.GetOwnerIdByHotelId(ctx, hotelReq)
		if err != nil {
			st, ok := status.FromError(err)
			if ok && st.Code() == codes.NotFound {
				b.log.Warn("error in service gRPC GetBookingsByHotelID:", zap.Error(myerror.ErrHotelNotFound))
				return fmt.Errorf("error in service GetOwnerIdByHotelId: %w", myerror.ErrHotelNotFound)
			}
			b.log.Error("error in service gRPC GetBookingsByHotelID:", zap.Error(err))
			return fmt.Errorf("error in service GetOwnerIdByHotelId: %w", err)
		}
		authReq := &authpb.GetHotelierRequest{OwnerID: hotelResponse.OwnerId}
		authResponse, err := b.authSvcClient.GetHotelierInformation(ctx, authReq)
		if err != nil {
			st, ok := status.FromError(err)
			if ok && st.Code() == codes.NotFound {
				b.log.Warn("error in service gRPC GetBookingsByHotelID:", zap.Error(err))
				return fmt.Errorf("error in service GetOwnerIdByHotelId: %w", myerror.ErrHotelNotFound)
			}
			b.log.Error("error in service gRPC GetBookingsByHotelID:", zap.Error(err))
			return fmt.Errorf("error in service GetOwnerIdByHotelId: %w", err)
		}

		hotelierMessage := bookingMessage.ToHotelierMessage(authResponse.Username, authResponse.ChatID)

		kafkaHotelierMessage, err := json.Marshal(hotelierMessage)
		if err != nil {
			b.log.Error("error in service UpdateBookingStatus", zap.Error(err))
			return fmt.Errorf("error in Marshal KafkaHotelierMessage: %w", err)
		}
		err = b.messageProducer.SendHotelierMessage(ctx, kafkaHotelierMessage)
		if err != nil {
			b.log.Error("error in service UpdateBookingStatus", zap.Error(err))
			return fmt.Errorf("error in SendMessage: %w", err)
		}
	case "fail":
		b.log.Warn("The payment is failing")
	}
	b.log.Info("in service update booking status end successfully", zap.String("status", BookingStatus))
	return nil
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
