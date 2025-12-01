package payment

import "context"

type TokenService interface {
	ValidateToken(ctx context.Context, token string) (*Payment, error)
	GenerateToken(ctx context.Context, payment *Payment) (string, error)
}
