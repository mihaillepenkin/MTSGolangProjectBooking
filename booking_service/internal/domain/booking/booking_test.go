package booking

import (
	"errors"
	"testing"

	error2 "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking/error"
	"gotest.tools/v3/assert"
)

func TestValidateAndGetBookingStatus(t *testing.T) {
	tests := []struct {
		status      string
		expectedErr error
	}{
		{"unpaid", nil},
		{"paid", nil},
		{"none", error2.ErrBookingStatusIsInCorrect},
	}

	for _, test := range tests {
		_, err := ValidateAndGetBookingStatus(test.status)

		assert.Assert(t, errors.Is(err, test.expectedErr))
	}
}
