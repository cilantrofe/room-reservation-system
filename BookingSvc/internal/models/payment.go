package models

type PaymentRequest struct {
	CardNumber string `json:"card_number"`
	Amount     int    `json:"amount"`
	WebHookURL string `json:"web_hook_url"`

	MetaData *BookingMessage `json:"meta_data"` //Это в meta data
}

type PaymentResponse struct {
	Status string `json:"status"`

	MetaData *BookingMessage `json:"meta_data"` //Это в meta data
}
