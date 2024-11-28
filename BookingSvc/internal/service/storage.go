package service

import (
	"context"
	"github.com/Quizert/room-reservation-system/BookingSvc/internal/models"
	"time"
)

type Storage interface {
	GetUnavailableRoomsByHotelId(ctx context.Context, HotelID int, startDate, endDate time.Time) (map[int]struct{}, error)
	GetBookingsByUserID(ctx context.Context, userID int) ([]*models.Booking, error)
	GetBookingsByHotelID(ctx context.Context, bookingID int) (*models.Booking, error)
	CreateBooking(ctx context.Context, booking *models.Booking) error
	UpdateBooking(ctx context.Context, booking *models.Booking) error
	DeleteBooking(ctx context.Context, bookingID int) error
}
