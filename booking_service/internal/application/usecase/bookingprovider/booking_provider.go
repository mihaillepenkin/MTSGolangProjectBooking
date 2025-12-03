package bookingprovider

import (
	"context"
	"log/slog"
	"time"

	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking"
	error2 "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking/error"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/hotel"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/user"
)

type BookingProvider struct {
	bookingRepo booking.Repository
	hotelRepo   hotel.Repository
	logger      *slog.Logger
}

func NewBookingProvider(bookingRepo booking.Repository, hotelRepo hotel.Repository) *BookingProvider {
	return &BookingProvider{bookingRepo: bookingRepo, hotelRepo: hotelRepo, logger: slog.Default().With("component", "booking_provider")}
}

func (b *BookingProvider) GetOccupiedRoomDurations(ctx context.Context, hotelName string, roomNumber string) ([][]time.Time, error) {
	durations, err := b.bookingRepo.GetDurationsByRoom(ctx, hotelName, roomNumber)
	if err != nil {
		b.logger.Error("Failed to get occupied durations", "error", err)
		return nil, err
	}

	return durations, nil
}

func (b *BookingProvider) GetBookingsByUser(ctx context.Context, user *user.User) ([]*booking.Booking, error) {
	bookings, err := b.bookingRepo.GetByUser(ctx, user.ID)

	if err != nil {
		b.logger.Error("Error while getting bookings by user: ", "error", err)
		return nil, err
	}

	return bookings, nil
}

func (b *BookingProvider) GetBookingsByHotelier(ctx context.Context, user *user.User, hotelName string) ([]*booking.Booking, error) {
	ok, err := b.hotelRepo.IsHotelier(ctx, user.ID, hotelName)
	if err != nil {
		b.logger.Error("Error while checking hotelier: ", "error", err)
		return nil, err
	}

	if !ok {
		b.logger.Debug("Hotelier is not valid", "hotelierID", user.ID, "hotelName", hotelName)
		return nil, error2.ErrHotelierIsNotValid
	}

	bookings, err := b.bookingRepo.GetByHotel(ctx, hotelName)
	if err != nil {
		b.logger.Error("Error while getting bookings by hotelier: ", "error", err)
		return nil, err
	}

	return bookings, nil
}
