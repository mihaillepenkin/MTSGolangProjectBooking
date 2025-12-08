package error

import "errors"

var (
	ErrHotelRoomIsNotFound      = errors.New("grpc room is not found")
	ErrHotelRoomInvalidArgument = errors.New("grpc room data is invalid")
	ErrHotelierInvalidArgument  = errors.New("hotelier data is not valid")
)
