package booking

import (
	"time"

	error2 "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking/error"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking/object"
)

type BookingStatus string

const (
	BookingStatusUnpaid BookingStatus = "unpaid"
	BookingStatusPaid   BookingStatus = "paid"
	BookingStatusNone   BookingStatus = ""
)

func ValidateAndGetBookingStatus(status string) (BookingStatus, error) {
	switch status {
	case string(BookingStatusUnpaid):
		return BookingStatusUnpaid, nil
	case string(BookingStatusPaid):
		return BookingStatusPaid, nil
	}
	return BookingStatusNone, error2.ErrBookingStatusIsInCorrect
}

type Booking struct {
	ID         object.BookingID
	UserID     string
	HotelName  string
	RoomNumber string
	TotalPrice float64
	Currency   string
	CheckIn    time.Time
	CheckOut   time.Time
	Status     BookingStatus
	PaymentID  string
}
