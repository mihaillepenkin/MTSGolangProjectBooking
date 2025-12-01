package payment

import "context"

type Service interface {
	ProcessPayment(ctx context.Context, token string) error
	CreatePayment(ctx context.Context, payment *Payment) (string, error)
}
