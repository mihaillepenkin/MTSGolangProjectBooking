package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/adapter/userkey"
	userdomain "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/user"
)

type AuthMiddleware struct {
	tokenService userdomain.TokenService
}

func NewAuthMiddleware(tokenService userdomain.TokenService) *AuthMiddleware {
	return &AuthMiddleware{tokenService: tokenService}
}

func (am *AuthMiddleware) Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			next.ServeHTTP(w, r)
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}
		token := parts[1]

		if token != "" {
			user, err := am.tokenService.ValidateToken(r.Context(), token)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), userkey.UserKey{}, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
