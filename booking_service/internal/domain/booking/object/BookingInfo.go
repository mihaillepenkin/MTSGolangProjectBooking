package object

import "time"

type BookingInfo struct {
	UserID     string
	HotelName  string
	RoomNumber string
	CheckIn    time.Time
	CheckOut   time.Time
}
