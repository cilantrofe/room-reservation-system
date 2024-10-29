package service

import (
	"BookingSvc/internal/models"
	"context"
)

type BookingService struct {
	service Repository
}

func NewBookingService(db Repository) *BookingService {
	return &BookingService{db}
}

func (b *BookingService) GetBookingsByUserID(ctx context.Context, userID int) ([]*models.Booking, error) {
	return b.service.GetBookingsByUserID(ctx, userID)
}

func (b *BookingService) GetBooking(ctx context.Context, id int) (*models.Booking, error) {
	return b.service.GetBooking(ctx, id)
}

func (b *BookingService) CreateBooking(ctx context.Context, booking *models.Booking) error {
	return b.service.CreateBooking(ctx, booking)
}

func (b *BookingService) UpdateBooking(ctx context.Context, booking *models.Booking) error {
	return b.service.UpdateBooking(ctx, booking)
}
func (b *BookingService) DeleteBooking(ctx context.Context, id int) error {
	return b.service.DeleteBooking(ctx, id)
}
