package models

type PaymentRequest struct {
	CardNumber string `json:"card_number"`
	Amount     int    `json:"amount"`
	BookingID  int    `json:"booking_id"`
	WebHookURL string `json:"web_hook_url"`
}

type PaymentResponse struct {
	BookingID int    `json:"booking_id"`
	Status    string `json:"status"`
}
