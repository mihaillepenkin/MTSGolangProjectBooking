package response

import "time"

type BookingResponse struct {
	UserID     string    `json:"userID"`
	HotelName  string    `json:"hotelName"`
	RoomNumber string    `json:"roomNumber"`
	TotalPrice float64   `json:"totalPrice"`
	Currency   string    `json:"currency"`
	CheckIn    time.Time `json:"checkIn"`
	CheckOut   time.Time `json:"checkOut"`
}
