package bookingsaver

import (
	"context"
	"log/slog"

	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/application/usecase/transactionmanager"
	bookingdomain "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking"
	error2 "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking/error"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking/object"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/hotel"
)

type BookingProvider struct {
	bookingRepo bookingdomain.Repository
	txManager   transactionmanager.TransactionManager[string]
	hotelRepo   hotel.Repository
}

func (b *BookingProvider) BookRoom(ctx context.Context, bookingInfo *object.BookingInfo) (string, error) {

	if bookingInfo.CheckOut.Before(bookingInfo.CheckIn) {
		return "", error2.ErrTimeDurationIsNotValid
	}

	return b.txManager.InTransaction(ctx, func(ctx context.Context) (string, error) {
		isIntersected, err := b.bookingRepo.IsIntersected(ctx, bookingInfo.HotelName, bookingInfo.RoomNumber, bookingInfo.CheckIn, bookingInfo.CheckOut)
		if err != nil {
			slog.Error("Error while checking intersected booking", "error", err)
			return "", err
		}

		if !isIntersected {
			return "", error2.ErrBookingIsIntersected
		}

		roomInfo, err := b.hotelRepo.GetRoomInfo(ctx, bookingInfo.HotelName, bookingInfo.RoomNumber)
		if err != nil {
			slog.Error("Error while getting room info", "error", err)
			return "", err
		}

		days := int64(bookingInfo.CheckOut.Sub(bookingInfo.CheckIn).Hours() / 24)

		booking := &bookingdomain.Booking{
			UserID:     bookingInfo.UserID,
			HotelName:  bookingInfo.HotelName,
			RoomNumber: bookingInfo.RoomNumber,
			TotalPrice: float64(days * roomInfo.Amount),
			Currency:   roomInfo.Currency,
			CheckIn:    bookingInfo.CheckIn,
			CheckOut:   bookingInfo.CheckOut,
			Status:     bookingdomain.BookingStatusUnpaid,
		}

		//TODO need to begin transaction from payment system and return url to user

		err = b.bookingRepo.Save(ctx, booking)
		if err != nil {
			slog.Error("Error while saving booking", "error", err)
			return "", err
		}

		return "url", nil
	})
}

func (b *BookingProvider) ConfirmBooking(ctx context.Context, bookingInfo *object.BookingInfo) error {
	_, err := b.txManager.InTransaction(ctx, func(ctx context.Context) (string, error) {
		booking, err := b.bookingRepo.GetByBookingInfo(ctx, bookingInfo)
		if err != nil {
			slog.Error("Error while getting booking", "error", err)
			return "", err
		}

		booking.Status = bookingdomain.BookingStatusPaid
		err = b.bookingRepo.Save(ctx, booking)
		if err != nil {
			slog.Error("Error while saving booking", "error", err)
			return "", err
		}

		return "", nil
	})

	if err != nil {
		slog.Error("Error while saving booking", "error", err)
		return err
	}

	return nil
}
