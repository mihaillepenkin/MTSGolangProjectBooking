package user

import "context"

type TokenService interface {
	ValidateToken(ctx context.Context, token string) (*User, error)
}
