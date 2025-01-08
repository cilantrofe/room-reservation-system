package handler

import (
	"github.com/Quizert/room-reservation-system/HotelSvc/internal/service"
	"github.com/Quizert/room-reservation-system/Libs/middleware"
	"net/http"
)

func RegisterHotelRoutes(mux *http.ServeMux, hotelService *service.HotelService, roomService *service.RoomService) {
	handler := &HotelHandler{hotelService: hotelService, roomService: roomService}

	middlewareHandler := middleware.NewMiddleware("LUIGI")
	mux.HandleFunc("/hotels", handler.GetHotels)                                       // GET - список отелей
	mux.HandleFunc("/add_hotel", middlewareHandler.Auth(handler.AddHotel, true))       // POST - добавление отеля
	mux.HandleFunc("/update_hotel", middlewareHandler.Auth(handler.UpdateHotel, true)) // PUT - обновление отеля
	mux.HandleFunc("/add_room", middlewareHandler.Auth(handler.AddRoom, true)) // POST - добавление комнаты в отель
	mux.HandleFunc("/add_room_type", middlewareHandler.Auth(handler.AddRoomType, true)) // POST - добавление типа комнаты
}
