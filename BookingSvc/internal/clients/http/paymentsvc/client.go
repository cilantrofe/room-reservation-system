package paymentsvc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Quizert/room-reservation-system/BookingSvc/internal/models"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Client struct {
	baseUrl string
	client  *http.Client
}

type Request struct {
	CardNumber string `json:"card_number"`
	Amount     int    `json:"amount"`
	WebHookURL string `json:"web_hook_url"`

	MetaData *models.BookingMessage `json:"meta_data"` //Это в meta data
}

type Response struct {
	Status string `json:"status"`

	MetaData *models.BookingMessage `json:"meta_data"` //Это в meta data
}

func ToPaymentRequest(bookingMessage *models.BookingMessage, cardNumber string, amount int) *Request {
	return &Request{
		CardNumber: cardNumber,
		Amount:     amount,
		WebHookURL: "http://booking-service:8080/bookings/payment/response?booking_id=" + strconv.Itoa(bookingMessage.BookingID),

		MetaData: bookingMessage,
	}
}

func NewPaymentSvcClient(baseUrl string) *Client {
	return &Client{
		baseUrl: baseUrl,
		client:  &http.Client{Timeout: 5 * time.Minute},
	}
}

func (c *Client) CreatePaymentRequest(ctx context.Context, paymentRequest *Request) error {
	jsonRequest, err := json.Marshal(paymentRequest)
	if err != nil {
		return fmt.Errorf("err in marshaling json: %w", err)
	}
	log.Println("JSON Request: ", string(jsonRequest))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseUrl, bytes.NewBuffer(jsonRequest))
	if err != nil {
		return fmt.Errorf("err in creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("err in sending request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("err in payment service status: %s", resp.Status)
	}
	return nil
}
