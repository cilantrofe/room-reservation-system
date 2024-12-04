package handlers

import (
	"encoding/json"
	"github.com/Quizert/room-reservation-system/PaymentSystem/internal/models"
	"github.com/Quizert/room-reservation-system/PaymentSystem/internal/service"
	"net/http"
)

type PaymentHandler struct {
	PaymentService *service.PaymentService
}

func NewPaymentHandler(PaymentService *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		PaymentService: PaymentService,
	}
}

func (p *PaymentHandler) ProcessPayment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var paymentRequest *models.PaymentRequest
	err := json.NewDecoder(r.Body).Decode(&paymentRequest)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	p.PaymentService.ProcessPayment(ctx, paymentRequest)

}
