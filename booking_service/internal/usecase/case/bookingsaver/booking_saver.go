package bookingsaver

import (
	"context"
	"log/slog"

	bookingdomain "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking"
	error2 "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking/error"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking/object"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/hotel"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/payment"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/usecase/case/transactionmanager"
)

type BookingSaver struct {
	bookingRepo       bookingdomain.Repository
	txManager         transactionmanager.TransactionManager[string]
	hotelRepo         hotel.Repository
	paymentSender     payment.PaymentSender
	webhookHandlerURL string
	logger            *slog.Logger
}

func NewBookingSaver(bookingRepo bookingdomain.Repository, txManager transactionmanager.TransactionManager[string], hotelRepo hotel.Repository, paymentSender payment.PaymentSender, webhookHandlerURL string) *BookingSaver {
	return &BookingSaver{bookingRepo: bookingRepo, txManager: txManager, hotelRepo: hotelRepo, paymentSender: paymentSender, webhookHandlerURL: webhookHandlerURL, logger: slog.Default().With("component", "booking_saver")}
}

func (b *BookingSaver) BookRoom(ctx context.Context, bookingInfo *object.BookingInfo) (string, error) {

	if bookingInfo.CheckOut.Before(bookingInfo.CheckIn) {
		return "", error2.ErrTimeDurationIsNotValid
	}

	return b.txManager.InTransaction(ctx, func(ctx context.Context) (string, error) {
		isIntersected, err := b.bookingRepo.IsIntersected(ctx, bookingInfo.HotelID, bookingInfo.RoomNumber, bookingInfo.CheckIn, bookingInfo.CheckOut)
		if err != nil {
			b.logger.Error("Error while checking intersected booking", "error", err)
			return "", err
		}

		if isIntersected {
			return "", error2.ErrBookingIsIntersected
		}

		roomInfo, err := b.hotelRepo.GetRoomInfo(ctx, bookingInfo.HotelID, bookingInfo.RoomNumber)
		if err != nil {
			b.logger.Error("Error while getting room info", "error", err)
			return "", err
		}

		days := int64(bookingInfo.CheckOut.Sub(bookingInfo.CheckIn).Hours() / 24)

		booking := &bookingdomain.Booking{
			UserID:     bookingInfo.User.ID,
			HotelID:    bookingInfo.HotelID,
			RoomNumber: bookingInfo.RoomNumber,
			TotalPrice: float64(days * roomInfo.Amount),
			Currency:   roomInfo.Currency,
			CheckIn:    bookingInfo.CheckIn,
			CheckOut:   bookingInfo.CheckOut,
			Status:     bookingdomain.BookingStatusUnpaid,
		}

		paymentInfo := &payment.PaymentInfo{BookingInfo: *bookingInfo, Price: booking.TotalPrice, Currency: roomInfo.Currency, URL: b.webhookHandlerURL}
		response, err := b.paymentSender.SendPayment(ctx, paymentInfo)
		if err != nil {
			b.logger.Error("Error while sending payment", "error", err)
			return "", err
		}

		booking.PaymentID = response.PaymentID

		err = b.bookingRepo.Save(ctx, booking)
		if err != nil {
			b.logger.Error("Error while saving booking", "error", err)
			return "", err
		}

		return response.URL, nil
	})
}

func (b *BookingSaver) DeleteBooking(ctx context.Context, bookingInfo *object.BookingInfo) error {
	_, err := b.txManager.InTransaction(ctx, func(ctx context.Context) (string, error) {
		booking, err := b.bookingRepo.GetByBookingInfo(ctx, bookingInfo)
		if err != nil {
			b.logger.Error("Error while getting booking", "error", err)
			return "", err
		}

		err = b.bookingRepo.Delete(ctx, booking)
		if err != nil {
			b.logger.Error("Error while deleting booking", "error", err)
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
			b.logger.Error("Error while getting booking", "error", err)
			return "", err
		}

		booking.Status = bookingdomain.BookingStatusPaid
		err = b.bookingRepo.Save(ctx, booking)
		if err != nil {
			b.logger.Error("Error while saving booking", "error", err)
			return "", err
		}

		return "", nil
	})

	if err != nil {
		b.logger.Error("Error while saving booking", "error", err)
		return err
	}

	return nil
}
