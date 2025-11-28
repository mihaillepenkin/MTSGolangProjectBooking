package bookingprovider

import (
	"context"
	"log/slog"
	"time"

	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking"
	error2 "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking/error"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/hotel"
)

type BookingProvider struct {
	bookingRepo booking.Repository
	hotelRepo   hotel.Repository
}

func NewBookingProvider(bookingRepo booking.Repository) *BookingProvider {
	return &BookingProvider{bookingRepo: bookingRepo}
}

func (b *BookingProvider) GetOccupiedRoomDurations(ctx context.Context, hotelName string, roomNumber string) ([][]time.Time, error) {
	durations, err := b.bookingRepo.GetDurationsByRoom(ctx, hotelName, roomNumber)
	if err != nil {
		slog.Error("Failed to get occupied durations", "error", err)
		return nil, err
	}

	return durations, nil
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
	ok, err := b.hotelRepo.IsHotelier(ctx, hotelierID, hotelName)
	if err != nil {
		slog.Error("Error while checking hotelier: ", "error", err)
		return nil, err
	}

	if !ok {
		slog.Debug("Hotelier is not valid", "hotelierID", hotelierID, "hotelName", hotelName)
		return nil, error2.ErrHotelierIsNotValid
	}

	bookings, err := b.bookingRepo.GetByHotel(ctx, hotelName)
	if err != nil {
		slog.Error("Error while getting bookings by hotelier: ", "error", err)
		return nil, err
	}

	return bookings, nil
}
