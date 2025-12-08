package jwtclaims

import "github.com/golang-jwt/jwt/v5"

type JWTClaims struct {
	UserID string `json:"sub"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}
