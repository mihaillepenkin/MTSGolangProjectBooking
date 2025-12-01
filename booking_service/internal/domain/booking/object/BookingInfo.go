package object

import (
	"time"

	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/user"
)

type BookingInfo struct {
	User       user.User `json:"user"`
	HotelName  string    `json:"hotel_name"`
	RoomNumber string    `json:"room_number"`
	CheckIn    time.Time `json:"check_in"`
	CheckOut   time.Time `json:"check_out"`
}
