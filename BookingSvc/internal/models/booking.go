package models

import "time"

type Booking struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	RoomID    int       `json:"room_id"`
	HotelID   int       `json:"hotel_id"`
	Status    string    `json:"status"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}
