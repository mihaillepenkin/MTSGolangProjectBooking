package request

import "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking/object"

type WebhookRequest struct {
	Price    float64            `json:"price"`
	Currency string             `json:"currency"`
	Info     object.BookingInfo `json:"metadata"`
	Status   string             `json:"status"`
}
