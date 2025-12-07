package response

import "time"

type BookingResponse struct {
	UserID     string    `json:"userID"`
	HotelID    int64     `json:"hotelID"`
	RoomNumber int64     `json:"roomNumber"`
	TotalPrice float64   `json:"totalPrice"`
	Currency   string    `json:"currency"`
	CheckIn    time.Time `json:"checkIn"`
	CheckOut   time.Time `json:"checkOut"`
}
