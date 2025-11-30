package booking

import (
	"context"
	"time"

	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/user"
)

type Provider interface {
	GetBookingsByUser(ctx context.Context, user *user.User) ([]*Booking, error)
	GetBookingsByHotelier(ctx context.Context, user *user.User, hotelName string) ([]*Booking, error)
	GetOccupiedRoomDurations(ctx context.Context, hotelName string, roomNumber string) ([][]time.Time, error)
}
