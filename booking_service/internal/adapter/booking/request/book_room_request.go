package request

import "time"

type BookRoomRequest struct {
	HotelID    int64     `json:"hotelID"`
	RoomNumber int64     `json:"roomNumber"`
	CheckIn    time.Time `json:"checkIn"`
	CheckOut   time.Time `json:"checkOut"`
}
