package handler

import (
	"BookingSvc/internal/service"
	"net/http"
)

func SetupRoutes(bookingService *service.BookingService) *http.ServeMux {
	mux := http.NewServeMux()

	bookingHandler := NewBookingHandler(bookingService)
	mux.HandleFunc("/bookings", bookingHandler.CreateBooking)
	mux.HandleFunc("/bookings/user", bookingHandler.GetBookingByUserId)
	return mux
}
