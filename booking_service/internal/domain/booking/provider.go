package booking

import (
	"context"
	"time"
)

type Provider interface {
	GetBookingsByUser(ctx context.Context, userID string) ([]*Booking, error)
	GetBookingsByHotelier(ctx context.Context, hotelierID string, hotelName string) ([]*Booking, error)
	GetOccupiedRoomDurations(ctx context.Context, hotelName string, roomNumber string) ([][]time.Time, error)
}
