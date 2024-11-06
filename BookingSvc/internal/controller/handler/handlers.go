package handler

import (
	"BookingSvc/internal/models"
	"BookingSvc/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type BookingHandler struct {
	bookingService *service.BookingService
}

func NewBookingHandler(b *service.BookingService) *BookingHandler {
	return &BookingHandler{b}
}

func (b *BookingHandler) GetBookingByUserID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		fmt.Println(userId)
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}
	bookings, err := b.bookingService.GetBookingsByUserID(ctx, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookings)
}

func (b *BookingHandler) GetBookingByHotelID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	hotelID, err := strconv.Atoi(r.URL.Query().Get("hotel_id"))
	if err != nil {
		fmt.Println(hotelID)
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}
	bookings, err := b.bookingService.GetBookingsByHotelID(ctx, hotelID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookings)
}

func (b *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	var booking models.Booking

	if err := json.NewDecoder(r.Body).Decode(&booking); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}
	fmt.Println(booking)
	ctx := r.Context()
	if err := b.bookingService.CreateBooking(ctx, &booking); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Booking created successfully",
		"id":      booking.ID,
	})
	if err != nil {
		return
	}
}
