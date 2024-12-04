package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Quizert/room-reservation-system/PaymentSystem/internal/models"
	"log"
	"net/http"
	"time"
)

type PaymentService struct {
}

func NewPaymentService() *PaymentService {
	return &PaymentService{}
}

func (p *PaymentService) sendWebHook(ctx context.Context, url string, paymentResponse *models.PaymentResponse) error {
	data, err := json.Marshal(paymentResponse)
	if err != nil {
		return fmt.Errorf("failed to marshal paymentResponse: %v", err)
	}
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(data))
	log.Println(url)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	log.Println(resp.Body, resp.StatusCode)
	defer resp.Body.Close()
	return nil
}

func (p *PaymentService) ProcessPayment(ctx context.Context, req *models.PaymentRequest) error {
	time.Sleep(15 * time.Second) // Имитация обратки платежа, связь с банком и т.д.
	paymentResponse := &models.PaymentResponse{
		BookingID: req.BookingID,
		Status:    "success",
	}
	err := p.sendWebHook(ctx, req.WebHookURL, paymentResponse)
	if err != nil {
		return err
	}
	return nil
}
