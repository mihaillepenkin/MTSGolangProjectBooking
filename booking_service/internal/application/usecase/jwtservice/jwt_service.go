package jwtservice

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/application/dto/jwtclaims"
	userdomain "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/user"
	error2 "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/user/error"
)

type JwtService struct {
	secretKey string
}

func NewJwtService(secretKey string) *JwtService {
	return &JwtService{secretKey: secretKey}
}

func (j *JwtService) ValidateToken(ctx context.Context, token string) (*userdomain.User, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &jwtclaims.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(j.secretKey), nil
	})

	if err != nil {
		slog.Error("validate token error", "error", err)
		return &userdomain.User{}, error2.ErrFailedToAuthorizeUser
	}

	if claims, ok := jwtToken.Claims.(*jwtclaims.JWTClaims); ok && jwtToken.Valid {
		user := &userdomain.User{ID: claims.UserID, Email: claims.Email, Role: claims.Role}
		err = userdomain.ValidateUser(user)
		if err != nil {
			slog.Error("user is incorrect", "error", err)
			return &userdomain.User{}, err
		}
		slog.Debug("validation of token is successful with ID", "ID", claims.UserID)
		return user, nil
	}

	slog.Error("validate token error", "error", err)
	return &userdomain.User{}, error2.ErrFailedToAuthorizeUser
}
