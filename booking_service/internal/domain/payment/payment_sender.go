package payment

import "context"

type PaymentSender interface {
	SendPayment(ctx context.Context, info PaymentInfo) (string, error)
}
