package handler

import (
	"github.com/Quizert/room-reservation-system/HotelSvc/internal/service"
	"net/http"
)

func RegisterHotelRoutes(mux *http.ServeMux, hotelService *service.HotelService) {
	handler := &HotelHandler{hotelService: hotelService}

	mux.HandleFunc("/hotels", handler.GetHotels)         // GET - список отелей
	mux.HandleFunc("/add_hotel", handler.AddHotel)       // POST - добавление отеля
	mux.HandleFunc("/update_hotel", handler.UpdateHotel) // PUT - обновление отеля
}
