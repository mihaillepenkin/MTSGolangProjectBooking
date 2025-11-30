package object

import (
	"time"

	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/user"
)

type BookingInfo struct {
	User       user.User
	HotelName  string
	RoomNumber string
	CheckIn    time.Time
	CheckOut   time.Time
}
