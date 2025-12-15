package object

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	error3 "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking/error"
	"gotest.tools/v3/assert"
)

func TestNewBookingID(t *testing.T) {
	id := uuid.New().String()

	_, err := NewBookingID(id)
	assert.NilError(t, err)

	_, err = NewBookingID("")
	assert.Assert(t, errors.Is(err, error3.ErrBookingIDCreatingIsNotValid))
}

func TestBookingID_ID(t *testing.T) {
	id := uuid.New().String()

	bookingID, _ := NewBookingID(id)

	assert.Equal(t, bookingID.ID(), id)
}

func TestBookingID_IsEmpty(t *testing.T) {
	id := uuid.New().String()

	bookingID, _ := NewBookingID(id)

	assert.Equal(t, bookingID.IsEmpty(), false)

	bookingID = BookingID{}

	assert.Equal(t, bookingID.IsEmpty(), true)
}
