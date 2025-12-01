package booking

import (
	"context"

	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking/object"
)

type Saver interface {
	BookRoom(ctx context.Context, bookingInfo *object.BookingInfo) (string, error)
	DeleteBooking(ctx context.Context, bookingInfo *object.BookingInfo) error
	ConfirmBooking(ctx context.Context, bookingInfo *object.BookingInfo) error
}
