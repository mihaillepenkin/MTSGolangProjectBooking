package error

import "errors"

var (
	ErrHotelRoomIsNotFound      = errors.New("hotel room is not found")
	ErrHotelRoomInvalidArgument = errors.New("hotel room data is invalid")
	ErrHotelierInvalidArgument  = errors.New("hotelier data is not valid")
)
