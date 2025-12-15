package jwtservice

import (
	"context"
	"errors"
	"log"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	error2 "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/user/error"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/usecase/dto/jwtclaims"
	"gotest.tools/v3/assert"
)

var (
	secretKey      = "secret"
	testJWTService *JwtService
)

func TestMain(m *testing.M) {
	testJWTService = NewJwtService(secretKey)

	log.Println("Running tests...")

	code := m.Run()

	os.Exit(code)
}

func TestJwtService_ShouldReturnFailedToAuthorizeErrValidateToken(t *testing.T) {
	ctx := context.Background()
	_, err := testJWTService.ValidateToken(ctx, "string")
	assert.Assert(t, errors.Is(error2.ErrFailedToAuthorizeUser, err))
}

func TestJwtService_ShouldReturnValidationErrValidateToken(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	duration := 2 * time.Hour
	claims := jwtclaims.JWTClaims{Name: "", Role: "", Email: "", UserID: "",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			IssuedAt:  jwt.NewNumericDate(now)},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		t.Fatal(err)
	}

	_, err = testJWTService.ValidateToken(ctx, tokenString)
	assert.Assert(t, err != nil)
}

func TestJwtService_ValidateToken(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	duration := 2 * time.Hour
	claims := jwtclaims.JWTClaims{Name: "1", Role: "1", Email: "123@mail.com", UserID: "1",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			IssuedAt:  jwt.NewNumericDate(now)},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		t.Fatal(err)
	}

	user, err := testJWTService.ValidateToken(ctx, tokenString)
	assert.NilError(t, err)

	assert.Assert(t, user.ID == claims.UserID && user.Email == claims.Email && user.Role == claims.Role && user.Name == claims.Name)
}
