package response

import (
	"testing"

	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking"
	"gotest.tools/v3/assert"
)

func TestNewBookingsResponse(t *testing.T) {
	bookings := []*booking.Booking{
		&booking.Booking{
			UserID: "1",
		},
	}

	response := NewBookingsResponse(bookings)

	assert.Assert(t, response.Bookings[0].UserID == bookings[0].UserID)
}
