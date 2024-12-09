package models

type BookingMessage struct {
	BookingID       int    `json:"booking_id"`
	HotelName       string `json:"hotel_name"`
	RoomDescription string `json:"room_description"`
	RoomNumber      int    `json:"room_number"`
	UserName        string `json:"user_name"`
	ChatId          string `json:"chat_id"`
}

func (req *BookingRequest) ToBookingMessage(bookingID int) *BookingMessage {
	return &BookingMessage{
		BookingID:       bookingID,
		HotelName:       req.HotelName,
		RoomDescription: req.RoomDescription,
		RoomNumber:      req.RoomNumber,
		UserName:        req.UserName,
		ChatId:          req.ChatId,
	}
}

func (message *BookingMessage) ToHotelierMessage(hotelierName string, hotelierChatID string) *BookingMessage {
	return &BookingMessage{
		BookingID:       message.BookingID,
		HotelName:       message.HotelName,
		RoomDescription: message.RoomDescription,
		RoomNumber:      message.RoomNumber,
		UserName:        hotelierName,
		ChatId:          hotelierChatID,
	}
}
