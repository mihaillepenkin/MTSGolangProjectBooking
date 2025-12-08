package error

import "errors"

var (
	ErrBookingStatusIsInCorrect    = errors.New("postgres status is in correct")
	ErrBookingIDCreatingIsNotValid = errors.New("postgres id is not valid")
	ErrBookingIsNotFound           = errors.New("postgres is not found")
	ErrBookingIsIntersected        = errors.New("postgres is intersected")
	ErrHotelierIsNotValid          = errors.New("hotelier is not valid")
	ErrTimeDurationIsNotValid      = errors.New("time duration is not valid")
)
