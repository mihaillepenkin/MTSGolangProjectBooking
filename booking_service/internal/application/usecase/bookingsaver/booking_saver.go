package bookingsaver

import (
	"context"
	"log/slog"

	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/application/usecase/transactionmanager"
	bookingdomain "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking"
	error2 "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking/error"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking/object"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/hotel"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/payment"
)

type BookingSaver struct {
	bookingRepo            bookingdomain.Repository
	txManager              transactionmanager.TransactionManager[string]
	hotelRepo              hotel.Repository
	paymentSender          payment.PaymentSender
	webhookHandlerEndpoint string
}

func (b *BookingSaver) BookRoom(ctx context.Context, bookingInfo *object.BookingInfo) (string, error) {

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
			UserID:     bookingInfo.User.ID,
			HotelName:  bookingInfo.HotelName,
			RoomNumber: bookingInfo.RoomNumber,
			TotalPrice: float64(days * roomInfo.Amount),
			Currency:   roomInfo.Currency,
			CheckIn:    bookingInfo.CheckIn,
			CheckOut:   bookingInfo.CheckOut,
			Status:     bookingdomain.BookingStatusUnpaid,
		}

		paymentInfo := payment.PaymentInfo{BookingInfo: *bookingInfo, Price: booking.TotalPrice, Currency: roomInfo.Currency, URL: b.webhookHandlerEndpoint}
		url, err := b.paymentSender.SendPayment(ctx, paymentInfo)
		if err != nil {
			slog.Error("Error while sending payment", "error", err)
			return "", err
		}

		err = b.bookingRepo.Save(ctx, booking)
		if err != nil {
			slog.Error("Error while saving booking", "error", err)
			return "", err
		}

		return url, nil
	})
}

func (b *BookingSaver) DeleteBooking(ctx context.Context, bookingInfo *object.BookingInfo) error {
	_, err := b.txManager.InTransaction(ctx, func(ctx context.Context) (string, error) {
		booking, err := b.bookingRepo.GetByBookingInfo(ctx, bookingInfo)
		if err != nil {
			slog.Error("Error while getting booking", "error", err)
			return "", err
		}

		err = b.bookingRepo.Delete(ctx, booking)
		if err != nil {
			slog.Error("Error while deleting booking", "error", err)
			return "", err
		}

		return "", nil
	})
	return err
}

func (b *BookingSaver) ConfirmBooking(ctx context.Context, bookingInfo *object.BookingInfo) error {
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
