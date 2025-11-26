package bookingprovider

import (
	"context"
	"log/slog"

	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking"
)

type BookingProvider struct {
	bookingRepo booking.Repository
}

func NewBookingProvider(bookingRepo booking.Repository) *BookingProvider {
	return &BookingProvider{bookingRepo: bookingRepo}
}

func (b *BookingProvider) GetBookingsByUser(ctx context.Context, userID string) ([]*booking.Booking, error) {
	bookings, err := b.bookingRepo.GetByUser(ctx, userID)

	if err != nil {
		slog.Error("Error while getting bookings by user: ", "error", err)
		return nil, err
	}

	return bookings, nil
}

func (b *BookingProvider) GetBookingsByHotelier(ctx context.Context, hotelierID string, hotelName string) ([]*booking.Booking, error) {
	//TODO need to add query to grpc server to check hotelier

	bookings, err := b.bookingRepo.GetByHotel(ctx, hotelName)
	if err != nil {
		slog.Error("Error while getting bookings by hotelier: ", "error", err)
		return nil, err
	}

	return bookings, nil
}
