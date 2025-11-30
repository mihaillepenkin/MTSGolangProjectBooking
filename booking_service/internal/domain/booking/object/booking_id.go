package object

import (
	error2 "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/error"

	"github.com/google/uuid"
)

type BookingID struct {
	id string
}

func NewBookingID(id string) (BookingID, error) {
	if _, err := uuid.Parse(id); err != nil {
		return BookingID{}, error2.ErrBookingIDCreatingIsNotValid
	}
	return BookingID{id: id}, nil
}

func (b BookingID) ID() string {
	return b.id
}

func (b BookingID) IsEmpty() bool {
	return b.id == ""
}
