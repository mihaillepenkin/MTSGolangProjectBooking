package user

import (
	"errors"
	"testing"

	error2 "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/user/error"
	"gotest.tools/v3/assert"
)

func TestValidateUser(t *testing.T) {
	user := &User{
		ID:   "",
		Role: "",
	}

	err := ValidateUser(user)
	assert.Assert(t, errors.Is(err, error2.ErrUserValidationFailed))

	user = &User{
		ID:    "1",
		Role:  "client",
		Email: "1",
	}

	err = ValidateUser(user)
	assert.Assert(t, errors.Is(err, error2.ErrUserValidationFailed))
	user = &User{
		ID:    "1",
		Role:  "client",
		Email: "123@mail.com",
	}

	err = ValidateUser(user)

	assert.Assert(t, err == nil)
}

func TestIsClient(t *testing.T) {
	user := &User{
		Role: "client",
	}

	assert.Equal(t, IsClient(user), true)
}

func TestIsHotelier(t *testing.T) {
	user := &User{
		Role: "hotelier",
	}

	assert.Equal(t, IsHotelier(user), true)
}
