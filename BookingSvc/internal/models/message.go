package models

type BookingMessage struct {
	BookingID       int    `json:"booking_id"`
	HotelName       string `json:"hotel_name"`
	RoomDescription string `json:"room_description"`
	RoomNumber      int    `json:"room_number"`
	Username        string `json:"user_name"`
	ChatID          string `json:"chat_id"`
}

func (req *BookingRequest) ToBookingMessage(bookingID int, username, chatID string) *BookingMessage {
	return &BookingMessage{
		BookingID:       bookingID,
		HotelName:       req.HotelName,
		RoomDescription: req.RoomDescription,
		RoomNumber:      req.RoomNumber,
		Username:        username,
		ChatID:          chatID,
	}
}

func (message *BookingMessage) ToHotelierMessage(hotelierName string, hotelierChatID string) *BookingMessage {
	return &BookingMessage{
		BookingID:       message.BookingID,
		HotelName:       message.HotelName,
		RoomDescription: message.RoomDescription,
		RoomNumber:      message.RoomNumber,
		Username:        hotelierName,
		ChatID:          hotelierChatID,
	}
}
