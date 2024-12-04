package paymentsvc

import (
	"bytes"
	"context"
	"encoding/json"
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
	BookingID  int    `json:"booking_id"`
	WebHookURL string `json:"web_hook_url"`
}

type Response struct {
	BookingID int    `json:"booking_id"`
	Status    string `json:"status"`
}

func NewPaymentSvcClient(baseUrl string) *Client {
	return &Client{
		baseUrl: baseUrl,
		client:  &http.Client{Timeout: 5 * time.Minute},
	}
}

func (c *Client) CreatePaymentRequest(ctx context.Context, cardNumber string, amount int, bookingID int) {
	if c.client == nil {
		log.Fatal("HTTP client is nil")
		return
	}
	log.Println(c.baseUrl, "DADADADADADADA", c.client, "CONTEXT", ctx)
	paymentRequest := Request{
		CardNumber: cardNumber,
		Amount:     amount,
		BookingID:  bookingID,
		WebHookURL: "http://booking-service:8080/bookings/payment/response?booking_id=" + strconv.Itoa(bookingID),
	}
	jsonRequest, err := json.Marshal(paymentRequest)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("JSON Request: ", string(jsonRequest))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseUrl, bytes.NewBuffer(jsonRequest))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
}
