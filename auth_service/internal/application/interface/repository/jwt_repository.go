package repository

import "time"

type JwtRepository interface {
	SaveRefreshToken(userID, token string, exp time.Time) error
	ValidateRefreshToken(token string) (string, error)
}