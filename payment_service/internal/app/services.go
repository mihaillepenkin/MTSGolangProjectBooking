package app

import (
	"net/http"
	"time"

	"github.com/mihaillepenkin/MTSGolangProjectBooking/payment_service/internal/application/usecase/jwt"
	payment2 "github.com/mihaillepenkin/MTSGolangProjectBooking/payment_service/internal/application/usecase/payment"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/payment_service/internal/application/usecase/webhooksender"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/payment_service/internal/domain/payment"
)

type Services struct {
	PaymentService payment.Service
	TokenService   payment.TokenService
	WebhookSender  payment.WebhookSender
}

func NewServices(secretKey string) *Services {
	tokenService := jwt.NewJWTService(secretKey)
	client := &http.Client{Timeout: 20 * time.Second}
	sender := webhooksender.NewWebhookSender(client)
	paymentService := payment2.NewPaymentService(tokenService, sender)
	return &Services{PaymentService: paymentService, TokenService: tokenService, WebhookSender: sender}
}
