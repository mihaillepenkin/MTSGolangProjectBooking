package booking

import (
	"context"
)

type Provider interface {
	GetBookingsByUser(ctx context.Context, userID string) ([]*Booking, error)
	GetBookingsByHotelier(ctx context.Context, hotelierID string, hotelName string) ([]*Booking, error)
}
