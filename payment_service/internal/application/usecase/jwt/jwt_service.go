package jwt

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/payment_service/internal/application/dto"
	paymentdomain "github.com/mihaillepenkin/MTSGolangProjectBooking/payment_service/internal/domain/payment"
	error2 "github.com/mihaillepenkin/MTSGolangProjectBooking/payment_service/internal/domain/payment/error"
)

type JWTService struct {
	secretKey string
}

func NewJWTService(secretKey string) *JWTService {
	return &JWTService{secretKey: secretKey}
}

func (j *JWTService) ValidateToken(ctx context.Context, token string) (*paymentdomain.Payment, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &dto.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(j.secretKey), nil
	})

	if err != nil {
		slog.Error("validate token err:", "error", err)
		return nil, err
	}

	if claims, ok := jwtToken.Claims.(*dto.JWTClaims); ok && jwtToken.Valid {
		payment := &paymentdomain.Payment{Price: claims.Payment.Price, Currency: claims.Payment.Currency, Metadata: claims.Payment.Metadata, URL: claims.Payment.URL}
		return payment, nil
	}

	slog.Error("validate token err")
	return nil, error2.ErrPaymentTokenIsInvalid
}

func (j *JWTService) GenerateToken(ctx context.Context, payment *paymentdomain.Payment) (string, error) {
	now := time.Now()
	duration := 5 * time.Minute
	claims := dto.JWTClaims{
		Payment: payment,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}
