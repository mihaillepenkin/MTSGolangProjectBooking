package postgres

import (
	"database/sql"
	"time"
)

type JWTRepoPG struct {
	DB *sql.DB
}

//now i don't realize this module

func (r *JWTRepoPG) SaveRefreshToken(userID, token string, exp time.Time) error {
	_, err := r.DB.Exec(
		"INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3) "+
			"ON CONFLICT (token) DO UPDATE SET expires_at = $3",
		userID, token, exp,
	)
	return err
}

func (r *JWTRepoPG) ValidateRefreshToken(token string) (string, error) {
	var userID string
	var expiresAt time.Time
	err := r.DB.QueryRow("SELECT user_id, expires_at FROM refresh_tokens WHERE token = $1", token).Scan(&userID, &expiresAt)
	if err != nil {
		return "", err
	}
	if time.Now().After(expiresAt) {
		return "", sql.ErrNoRows
	}
	return userID, nil
}