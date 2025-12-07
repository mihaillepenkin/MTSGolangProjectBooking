package booking

import (
	"context"
	"time"

	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/user"
)

type Provider interface {
	GetBookingsByUser(ctx context.Context, user *user.User) ([]*Booking, error)
	GetBookingsByHotelier(ctx context.Context, user *user.User, hotelID int64) ([]*Booking, error)
	GetOccupiedRoomDurations(ctx context.Context, hotelID int64, roomNumber int64) ([][]time.Time, error)
}
