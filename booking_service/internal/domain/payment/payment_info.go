package payment

import "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking/object"

type PaymentInfo struct {
	Price       float64            `json:"price"`
	Currency    string             `json:"currency"`
	BookingInfo object.BookingInfo `json:"metadata"`
	URL         string             `json:"url"`
}
