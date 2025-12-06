package booking

import (
	"context"
	"time"

	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking/object"
)

type Repository interface {
	Save(ctx context.Context, booking *Booking) error
	Delete(ctx context.Context, booking *Booking) error
	IsIntersected(ctx context.Context, hotelID int64, hotelRoom int64, checkIn time.Time, checkOut time.Time) (bool, error)
	GetByBookingInfo(ctx context.Context, bookingInfo *object.BookingInfo) (*Booking, error)
	GetByHotel(ctx context.Context, hotelID int64) ([]*Booking, error)
	GetByUser(ctx context.Context, userID string) ([]*Booking, error)
	GetDurationsByRoom(ctx context.Context, hotelID int64, roomNumber int64) ([][]time.Time, error)
	GetBookingsByStatus(ctx context.Context, status BookingStatus) ([]*Booking, error)
}
