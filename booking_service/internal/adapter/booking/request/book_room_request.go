package request

import "time"

type BookRoomRequest struct {
	HotelID  int64     `json:"hotelID"`
	RoomID   int64     `json:"roomID"`
	CheckIn  time.Time `json:"checkIn"`
	CheckOut time.Time `json:"checkOut"`
}
