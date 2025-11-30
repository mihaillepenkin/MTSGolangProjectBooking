package bookingsaver

import (
	"context"
	"log/slog"

	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/application/usecase/transactionmanager"
	bookingdomain "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking/object"
	error2 "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/error"
)

type BookingProvider struct {
	bookingRepo bookingdomain.Repository
	txManager   transactionmanager.TransactionManager[string]
}

func (b *BookingProvider) BookRoom(ctx context.Context, bookingInfo *object.BookingInfo) (string, error) {
	return b.txManager.InTransaction(ctx, func(ctx context.Context) (string, error) {
		isIntersected, err := b.bookingRepo.IsIntersected(ctx, bookingInfo.HotelName, bookingInfo.RoomNumber, bookingInfo.CheckIn, bookingInfo.CheckOut)
		if err != nil {
			slog.Error("Error while checking intersected booking", "error", err)
			return "", err
		}

		if !isIntersected {
			return "", error2.ErrBookingIsIntersected
		}

		//TODO need to add query to grpc server to get room info

		//TODO need to begin transaction from payment system and return url to user

		booking := &bookingdomain.Booking{
			UserID:     bookingInfo.UserID,
			HotelName:  bookingInfo.HotelName,
			RoomNumber: bookingInfo.RoomNumber,
			TotalPrice: 0,
			Currency:   "RUB",
			CheckIn:    bookingInfo.CheckIn,
			CheckOut:   bookingInfo.CheckOut,
			Status:     bookingdomain.BookingStatusUnpaid,
		}

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
