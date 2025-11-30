package request

import "time"

type BookRoomRequest struct {
	HotelName  string    `json:"hotelName"`
	RoomNumber string    `json:"roomNumber"`
	CheckIn    time.Time `json:"checkIn"`
	CheckOut   time.Time `json:"checkOut"`
}
