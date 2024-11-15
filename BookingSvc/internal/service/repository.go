package service

import (
	"BookingSvc/internal/models"
	"context"
	hotelSvc "github.com/Quizert/room-reservation-system/HotelSvc/api/grpc/hotelpb"
	"time"
)

type Repository interface {
	GetUnavailableRoomsByHotelId(ctx context.Context, HotelID int, startDate, endDate time.Time) (*[]hotelSvc.Room, error)
	GetBookingsByUserID(ctx context.Context, userID int) ([]*models.Booking, error)
	GetBookingsByHotelID(ctx context.Context, bookingID int) (*models.Booking, error)
	CreateBooking(ctx context.Context, booking *models.Booking) error
	UpdateBooking(ctx context.Context, booking *models.Booking) error
	DeleteBooking(ctx context.Context, bookingID int) error
}
