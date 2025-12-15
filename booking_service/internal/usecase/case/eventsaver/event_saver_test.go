package eventsaver

import (
	"context"
	"errors"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking/object"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/message"
	userdomain "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/user"
	"gotest.tools/v3/assert"
)

var (
	testSaver       *EventSaver
	testBookingInfo = &object.BookingInfo{User: userdomain.User{ID: uuid.New().String(), Role: "1", Name: "1", Email: "123@mail.com"},
		HotelID: 1, RoomID: 1, CheckIn: time.Now().UTC(), CheckOut: time.Now().UTC(), HotelName: "1", RoomNumber: 1}
)

type mockSaver struct{}

func (m *mockSaver) BookRoom(ctx context.Context, bookingInfo *object.BookingInfo) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockSaver) DeleteBooking(ctx context.Context, bookingInfo *object.BookingInfo) error {
	if bookingInfo.User.ID != testBookingInfo.User.ID {
		return errors.New("booking info not found")
	}

	return nil
}

func (m *mockSaver) ConfirmBooking(ctx context.Context, bookingInfo *object.BookingInfo) error {
	if bookingInfo.User.ID != testBookingInfo.User.ID {
		return errors.New("booking info not found")
	}
	return nil
}

type mockProducer struct{}

func (m *mockProducer) SendMessage(ctx context.Context, message *message.Message) error {
	return nil
}

func TestMain(m *testing.M) {
	testSaver = NewEventSaver(&mockProducer{}, &mockSaver{})

	log.Println("Running tests...")
	code := m.Run()

	os.Exit(code)
}

func TestEventSaver_DeleteBooking(t *testing.T) {
	ctx := context.Background()

	bookingInfo := &object.BookingInfo{User: userdomain.User{ID: "1"}}
	err := testSaver.DeleteBooking(ctx, bookingInfo)

	assert.Assert(t, err != nil)

	err = testSaver.DeleteBooking(ctx, testBookingInfo)

	assert.Assert(t, err == nil)
}

func TestEventSaver_ConfirmBooking(t *testing.T) {
	ctx := context.Background()

	bookingInfo := &object.BookingInfo{User: userdomain.User{ID: "1"}}
	err := testSaver.ConfirmBooking(ctx, bookingInfo)

	assert.Assert(t, err != nil)

	err = testSaver.ConfirmBooking(ctx, testBookingInfo)

	assert.Assert(t, err == nil)
}
