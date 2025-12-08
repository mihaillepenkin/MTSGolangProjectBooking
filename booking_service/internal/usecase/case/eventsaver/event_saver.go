package eventsaver

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking/object"
	messagedomain "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/message"
)

type EventSaver struct {
	producer messagedomain.Producer
	booking.Saver
	logger *slog.Logger
}

func NewEventSaver(producer messagedomain.Producer, saver booking.Saver) *EventSaver {
	return &EventSaver{producer: producer, Saver: saver, logger: slog.Default().With("component", "event_saver")}
}
func (s *EventSaver) DeleteBooking(ctx context.Context, bookingInfo *object.BookingInfo) error {
	err := s.Saver.DeleteBooking(ctx, bookingInfo)
	if err != nil {
		return err
	}

	message := messagedomain.NewMessage(bookingInfo.User.Email, bookingInfo.User.Name, getOperationInfo(bookingInfo), `
Payment did not go through the system`, messagedomain.StatusError)
	err = s.producer.SendMessage(ctx, message)
	if err != nil {
		s.logger.Error("Error sending message", "error", err)
	}

	return nil
}

func (s *EventSaver) ConfirmBooking(ctx context.Context, bookingInfo *object.BookingInfo) error {
	err := s.Saver.ConfirmBooking(ctx, bookingInfo)
	if err != nil {
		return err
	}

	message := messagedomain.NewMessage(bookingInfo.User.Email, bookingInfo.User.Name, getOperationInfo(bookingInfo), "",
		messagedomain.StatusOK)
	err = s.producer.SendMessage(ctx, message)
	if err != nil {
		s.logger.Error("Error sending message", "error", err)
	}
	return nil
}

func getOperationInfo(info *object.BookingInfo) string {
	return fmt.Sprintf("%v %v %v %v", info.HotelName, info.RoomNumber, info.CheckOut.Format(time.DateOnly), info.CheckIn.Format(time.DateOnly))
}
