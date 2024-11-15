package handler

import (
	"BookingSvc/internal/service"
	"net/http"
)

func SetupRoutes(bookingService *service.BookingService) *http.ServeMux {
	mux := http.NewServeMux()
	bookingHandler := NewBookingHandler(bookingService)

	mux.HandleFunc("/bookings", bookingHandler.CreateBooking)                  // POST - Создается новое бронирование
	mux.HandleFunc("/bookings/users", bookingHandler.GetBookingByUserID)       // GET - получаем все бронирования пользователя
	mux.HandleFunc("/bookings/hotels", bookingHandler.GetBookingByHotelID)     // Get - получаем все бронирования отельера
	mux.HandleFunc("/bookings/hotels/rooms", bookingHandler.GetAvailableRooms) //Тут добавить сортировку по времени
	return mux
}
