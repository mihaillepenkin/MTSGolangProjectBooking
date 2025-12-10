package payment

import (
	"context"
	"log/slog"
	"math/rand"

	paymentdomain "github.com/mihaillepenkin/MTSGolangProjectBooking/payment_service/internal/domain/payment"
	error2 "github.com/mihaillepenkin/MTSGolangProjectBooking/payment_service/internal/domain/payment/error"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/payment_service/internal/domain/payment/object"
)

type PaymentService struct {
	tokenService paymentdomain.TokenService
	sender       paymentdomain.WebhookSender
}

func NewPaymentService(tokenService paymentdomain.TokenService, sender paymentdomain.WebhookSender) *PaymentService {
	return &PaymentService{tokenService: tokenService, sender: sender}
}

func (p *PaymentService) ProcessPayment(ctx context.Context, token string) error {
	if token == "" {
		slog.Error("Token is empty")
		return error2.ErrPaymentTokenIsInvalid
	}

	payment, err := p.tokenService.ValidateToken(ctx, token)
	if err != nil {
		slog.Error("Token is invalid")
		return error2.ErrPaymentTokenIsInvalid
	}

	response := &object.PaymentResponse{Price: payment.Price, Currency: payment.Currency, Metadata: payment.Metadata, Status: "Failed"}
	if rand.Float64() < 0.75 {
		response.Status = "Success"
		err = p.sender.Send(ctx, payment.URL, response)
		if err != nil {
			slog.Error("Error sending payment webhook")
		}
	} else {
		err = p.sender.Send(ctx, payment.URL, response)
		if err != nil {
			slog.Error("Error sending payment webhook")
		}
	}

	if response.Status != "Success" {
		return error2.ErrPaymentFailed
	}

	return nil
}

func (p *PaymentService) CreatePayment(ctx context.Context, payment *paymentdomain.Payment) (string, error) {
	err := paymentdomain.ValidatePayment(payment)
	if err != nil {
		slog.Error("Error validating payment", "error", err)
		return "", err
	}

	token, err := p.tokenService.GenerateToken(ctx, payment)
	if err != nil {
		slog.Error("Error generating token", "error", err)
		return "", err
	}

	return token, nil
}
