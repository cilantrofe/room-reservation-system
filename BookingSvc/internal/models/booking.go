package models

import "time"

type BookingRequest struct {
	UserID    int       `json:"user_id"`
	RoomID    int       `json:"room_id"`
	HotelID   int       `json:"hotel_id"`
	Status    string    `json:"status"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`

	HotelName       string `json:"hotel_name"`
	RoomDescription string `json:"room_description"`
	RoomNumber      int    `json:"room_number"`
	UserName        string `json:"user_name"`
	ChatId          int    `json:"chat_id"`

	CardNumber string `json:"card_number"`
	Amount     int    `json:"amount"`
}

type BookingMessage struct {
	HotelName       string `json:"hotel_name"`
	RoomDescription string `json:"room_description"`
	RoomNumber      int    `json:"room_number"`
	UserName        string `json:"user_name"`
	ChatId          int    `json:"chat_id"`
}

type Booking struct {
	UserID    int       `json:"user_id"`
	RoomID    int       `json:"room_id"`
	HotelID   int       `json:"hotel_id"`
	Status    string    `json:"status"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

func (req *BookingRequest) ToBooking() *Booking {
	return &Booking{
		UserID:    req.UserID,
		RoomID:    req.RoomID,
		HotelID:   req.HotelID,
		Status:    req.Status,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}
}

func (req *BookingRequest) ToBookingMessage() *BookingMessage {
	return &BookingMessage{
		HotelName:       req.HotelName,
		RoomDescription: req.RoomDescription,
		RoomNumber:      req.RoomNumber,
		UserName:        req.UserName,
		ChatId:          req.ChatId,
	}
}
