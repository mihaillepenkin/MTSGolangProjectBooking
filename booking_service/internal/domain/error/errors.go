package error

import "errors"

var (
	ErrBookingStatusIsInCorrect    = errors.New("booking status is in correct")
	ErrBookingIDCreatingIsNotValid = errors.New("booking id is not valid")
	ErrBookingIsNotFound           = errors.New("booking is not found")
	ErrHotelierIsNotCorrect        = errors.New("hotelier is not correct")
	ErrBookingIsIntersected        = errors.New("booking is intersected")
)
