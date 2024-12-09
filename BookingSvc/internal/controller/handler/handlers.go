package handler

import (
	"encoding/json"
	"fmt"
	paymentClient "github.com/Quizert/room-reservation-system/BookingSvc/internal/clients/http/paymentsvc"
	"github.com/Quizert/room-reservation-system/BookingSvc/internal/models"
	"github.com/Quizert/room-reservation-system/BookingSvc/internal/service"
	"log"
	"net/http"
	"strconv"
	"time"
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
	//TODO: Если пользователь зареган в системе, то возвращаем 200, иначе ошибка 401 или 403 или 404
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}
	bookings, err := b.bookingService.GetBookingsByUserID(ctx, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookings)
}

func (b *BookingHandler) GetAvailableRooms(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	hotelId, err := strconv.Atoi(r.URL.Query().Get("hotel_id"))
	if err != nil {
		http.Error(w, "Invalid hotel_id", http.StatusBadRequest)
		return
	}
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	// Парсим start_date
	startDate, err := time.Parse(time.RFC3339, startDateStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid start_date: %v", err), http.StatusBadRequest)
		return
	}
	// Парсим end_date
	endDate, err := time.Parse(time.RFC3339, endDateStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid end_date: %v", err), http.StatusBadRequest)
		return
	}
	availableRooms, err := b.bookingService.GetAvailableRooms(ctx, hotelId, startDate.UTC(), endDate.UTC())
	if err != nil {
		// тут мб лог
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(availableRooms)

}

func (b *BookingHandler) GetBookingByHotelID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	hotelID, err := strconv.Atoi(r.URL.Query().Get("hotel_id"))
	if err != nil {
		fmt.Println(hotelID)
		http.Error(w, "Invalid user id", http.StatusBadRequest)
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
	ctx := r.Context()
	var bookingRequest models.BookingRequest
	if err := json.NewDecoder(r.Body).Decode(&bookingRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) //err.Error() - исправить
		return
	}

	if err := b.bookingService.CreateBooking(ctx, &bookingRequest); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) //Тут добавить обработку: бронирвание уже существует error
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (b *BookingHandler) HandlePaymentWebHook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var paymentResponse paymentClient.Response
	if err := json.NewDecoder(r.Body).Decode(&paymentResponse); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest) //err.Error() - исправить
		return
	}
	log.Println(paymentResponse)
	err := b.bookingService.UpdateBookingStatus(ctx, paymentResponse.Status, paymentResponse.MetaData)

	switch paymentResponse.Status {
	case "success":
		if err != nil {
			log.Println("handler UpdateBookingStatusSuccess: ", err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError) //err.Error() - исправить
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success booking!"))
	}
}
