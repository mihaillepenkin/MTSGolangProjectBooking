package dto

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/payment_service/internal/domain/payment"
)

type JWTClaims struct {
	Payment *payment.Payment `json:"payment"`
	jwt.RegisteredClaims
}
