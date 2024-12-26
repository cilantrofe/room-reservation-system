package handler

import (
	"encoding/json"
	"github.com/Quizert/room-reservation-system/HotelSvc/internal/models"
	"github.com/Quizert/room-reservation-system/HotelSvc/internal/service"
	"net/http"
)

type HotelHandler struct {
	hotelService *service.HotelService
	roomService  *service.RoomService
}

// GetHotels - обработчик для получения списка отелей
func (h *HotelHandler) GetHotels(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		hotels, err := h.hotelService.GetAllHotels()
		if err != nil {
			http.Error(w, "Failed to get hotels", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(hotels)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// AddHotel - обработчик для добавления нового отеля
func (h *HotelHandler) AddHotel(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		ctx := r.Context()
		ownerID := ctx.Value("user_id").(int)

		var hotel models.Hotel
		if err := json.NewDecoder(r.Body).Decode(&hotel); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		hotel.OwnerId = ownerID
		if err := h.hotelService.AddHotel(hotel); err != nil {
			http.Error(w, "Failed to add hotel", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// UpdateHotel - обработчик для обновления информации об отеле
func (h *HotelHandler) UpdateHotel(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		var hotel models.Hotel
		if err := json.NewDecoder(r.Body).Decode(&hotel); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		if err := h.hotelService.UpdateHotel(r.Context(), hotel); err != nil {
			http.Error(w, "Failed to update hotel", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *HotelHandler) AddRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var room models.Room
		if err := json.NewDecoder(r.Body).Decode(&room); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		if err := h.roomService.AddRoom(r.Context(), room); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *HotelHandler) AddRoomType(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var roomType models.RoomType
		if err := json.NewDecoder(r.Body).Decode(&roomType); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		if err := h.roomService.AddRoomType(roomType); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
