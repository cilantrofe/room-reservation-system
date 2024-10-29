package service

import (
	"BookingSvc/internal/models"
	"context"
)

type Repository interface {
	GetBooking(ctx context.Context, bookingID int) (*models.Booking, error)
	GetBookingsByUserID(ctx context.Context, userID int) ([]*models.Booking, error)
	CreateBooking(ctx context.Context, booking *models.Booking) error
	UpdateBooking(ctx context.Context, booking *models.Booking) error
	DeleteBooking(ctx context.Context, bookingID int) error
}
