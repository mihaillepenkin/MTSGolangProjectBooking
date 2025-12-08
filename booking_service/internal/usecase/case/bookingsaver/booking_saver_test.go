package bookingsaver

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking"
	error3 "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking/error"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking/object"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/hotel"
	error2 "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/hotel/error"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/payment"
	userdomain "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/user"
	"gotest.tools/v3/assert"
)

var (
	testWebhookHandler       = "webhookHandler"
	testDurations            = [][]time.Time{{time.Date(2025, 10, 10, 0, 0, 0, 0, time.UTC), time.Date(2025, 10, 11, 0, 0, 0, 0, time.UTC)}}
	testHotelName      int64 = 1
	testUserID               = "1"
	testRoomNumber     int64 = 1
	testUser                 = &userdomain.User{ID: testUserID}
	testBooking              = &booking.Booking{UserID: testUserID, HotelID: testHotelName, RoomID: testRoomNumber, TotalPrice: 100, Currency: "USD", CheckIn: testDurations[0][0], CheckOut: testDurations[0][1], Status: booking.BookingStatusPaid, PaymentID: "1"}
	testBookingInfo          = &object.BookingInfo{User: *testUser, HotelID: testHotelName, RoomID: testRoomNumber, CheckIn: testDurations[0][0], CheckOut: testDurations[0][1]}
	testSaver          *BookingSaver
)

type MockBookRepo struct{}

func (m *MockBookRepo) Save(ctx context.Context, booking *booking.Booking) error {
	return nil
}

func (m *MockBookRepo) Delete(ctx context.Context, booking *booking.Booking) error {
	return nil
}

func (m *MockBookRepo) IsIntersected(ctx context.Context, hotelName int64, hotelRoom int64, checkIn time.Time, checkOut time.Time) (bool, error) {
	if hotelName != testHotelName {
		return true, nil
	}
	return false, nil
}

func (m *MockBookRepo) GetByBookingInfo(ctx context.Context, bookingInfo *object.BookingInfo) (*booking.Booking, error) {
	if bookingInfo.HotelID == testHotelName && bookingInfo.RoomID == testRoomNumber {
		return testBooking, nil
	}
	return nil, fmt.Errorf("postgres info not found")
}

func (m *MockBookRepo) GetByHotel(ctx context.Context, hotelName int64) ([]*booking.Booking, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockBookRepo) GetByUser(ctx context.Context, userID string) ([]*booking.Booking, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockBookRepo) GetDurationsByRoom(ctx context.Context, hotelName int64, roomNumber int64) ([][]time.Time, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockBookRepo) GetBookingsByStatus(ctx context.Context, status booking.BookingStatus) ([]*booking.Booking, error) {
	//TODO implement me
	panic("implement me")
}

type MockHotelRepo struct{}

func (m *MockHotelRepo) IsHotelier(ctx context.Context, userID string, hotelName int64) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockHotelRepo) GetRoomInfo(ctx context.Context, hotelName int64, roomNumber int64) (*hotel.RoomInfo, error) {
	if hotelName == testHotelName && roomNumber == testRoomNumber {
		return &hotel.RoomInfo{}, nil
	}
	return nil, error2.ErrHotelRoomIsNotFound
}

type MockPaymentSender struct{}

func (m *MockPaymentSender) SendPayment(ctx context.Context, info *payment.PaymentInfo) (*payment.PaymentResponse, error) {
	return &payment.PaymentResponse{}, nil
}

type MockTransactionManager[T any] struct{}

func (m *MockTransactionManager[T]) InTransaction(ctx context.Context, fn func(ctx context.Context) (T, error)) (T, error) {
	return fn(ctx)
}

func TestMain(m *testing.M) {
	testSaver = NewBookingSaver(&MockBookRepo{}, &MockTransactionManager[string]{}, &MockHotelRepo{},
		&MockPaymentSender{}, testWebhookHandler)
	log.Println("Running tests...")
	code := m.Run()
	os.Exit(code)
}

func TestBookingSaver_BookRoom(t *testing.T) {
	ctx := context.Background()
	_, err := testSaver.BookRoom(ctx, &object.BookingInfo{HotelID: 2})
	assert.Assert(t, errors.Is(err, error3.ErrBookingIsIntersected), "book room error")
	_, err = testSaver.BookRoom(ctx, &object.BookingInfo{HotelID: testHotelName, RoomID: 2})
	assert.Equal(t, error2.ErrHotelRoomIsNotFound, err, "error must be not found")
	_, err = testSaver.BookRoom(ctx, testBookingInfo)
	assert.Assert(t, err == nil, "error should be nil")
}

func TestBookingSaver_ConfirmBooking(t *testing.T) {
	ctx := context.Background()
	err := testSaver.ConfirmBooking(ctx, testBookingInfo)
	assert.Assert(t, err == nil, "error should be nil")
}

func TestBookingSaver_GetBookingInfo(t *testing.T) {
	ctx := context.Background()
	err := testSaver.DeleteBooking(ctx, testBookingInfo)
	assert.Assert(t, err == nil, "error should be nil")
}
