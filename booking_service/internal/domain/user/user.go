package user

import (
	"net/mail"

	error2 "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/user/error"
)

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func ValidateUser(user *User) error {
	if user.ID == "" || user.Role == "" {
		return error2.ErrUserValidationFailed
	}

	_, err := mail.ParseAddress(user.Email)
	if err != nil {
		return error2.ErrUserValidationFailed
	}

	return nil
}

func IsClient(user *User) bool {
	return user.Role == "client"
}

func IsHotelier(user *User) bool {
	return user.Role == "hotelier"
}
