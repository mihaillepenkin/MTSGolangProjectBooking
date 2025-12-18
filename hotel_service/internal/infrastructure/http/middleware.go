package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserId string `json:"sub"`
	Role   string `json:"role"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization format", http.StatusBadRequest)
			return
		}

		token := parts[1]
		parsedToken, err := jwt.ParseWithClaims(token, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			http.Error(w, "token validation error", http.StatusUnauthorized)
			return
		}

		claims, ok := parsedToken.Claims.(*JWTClaims)
		if !ok || !parsedToken.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		if r.URL.Path == "/api/v1/hotels" {
			switch r.Method {
			case "GET":
				if claims.Role != "client" {
					http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
					return
				}
			case "POST":
				if claims.Role == "client" {
					http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
					return
				}
			case "PUT":
				if claims.Role == "client" {
					http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
					return
				}
			}
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "claims", claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
