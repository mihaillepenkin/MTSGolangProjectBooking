package object

import (
	"time"

	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/user"
)

type BookingInfo struct {
	User       user.User `json:"user"`
	HotelID    int64     `json:"hotel_id"`
	RoomID     int64     `json:"room_id"`
	CheckIn    time.Time `json:"check_in"`
	CheckOut   time.Time `json:"check_out"`
	HotelName  string    `json:"hotel_name"`
	RoomNumber int64     `json:"room_number"`
}
