package payment

import (
	"context"

	"github.com/mihaillepenkin/MTSGolangProjectBooking/payment_service/internal/domain/payment/object"
)

type WebhookSender interface {
	Send(ctx context.Context, URL string, response *object.PaymentResponse) error
}
