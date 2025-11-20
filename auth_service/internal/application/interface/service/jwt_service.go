package service

import (
	"auth_service/internal/domain/entity"
	"time"
)

type JwtService interface {
	GenerateAccessToken(user *entity.User) (string, time.Time, error)
	GenerateRefreshToken(user *entity.User) (string, time.Time, error)
	ValidateToken(token string) (*entity.User, error)
}