package usecase

import (
	"auth_service/internal/domain/entity"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTUseCase struct {
	secret string
}

func JWTService(secret string) *JWTUseCase {
	serv := new(JWTUseCase)
	serv.secret = secret
	return serv
}

//access token - токен на 24 часа для доступа к сервисам
func (j *JWTUseCase) GenerateAccessToken(user *entity.User) (string, time.Time, error) {
	exp := time.Now().Add(24 * time.Hour)
	claims := jwt.MapClaims{
		"sub": user.ID,
		"role": user.Role,
		"email": user.Email,
		"exp": exp.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(j.secret))
	return signed, exp, err
}

//refresh token - токен для получения access token
func (j *JWTUseCase) GenerateRefreshToken(user *entity.User) (string, time.Time, error) {
	exp := time.Now().Add(48 * time.Hour)
	claims := jwt.MapClaims{
		"sub": user.ID,
		"role": user.Role,
		"email": user.Email,
		"exp": exp.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(j.secret))
	return signed, exp, err
}

func (j *JWTUseCase) ValidateToken(tokenStr string) (*entity.User, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(j.secret), nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrInvalidKeyType
	}


	if exp, ok := claims["exp"].(float64); ok {
		if time.Unix(int64(exp), 0).Before(time.Now()) {
			return nil, fmt.Errorf("token expired")
		}
	}

	user := &entity.User{}

	if id, ok := claims["sub"].(string); ok {
		user.ID = id
	} else if idf, ok := claims["sub"].(float64); ok {
		user.ID = fmt.Sprintf("%d", int64(idf))
	}

	if value, ok := claims["role"].(string); ok {
		user.Role = value
	}

	if value, ok := claims["email"].(string); ok {
		user.Email = value
	}

	return user, nil
}