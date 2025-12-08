package bookingprovider

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking/object"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/hotel"
	userdomain "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/user"
	"gotest.tools/v3/assert"
)

var (
	testHotelName  int64 = 1
	testUserID           = "1"
	testRoomNumber int64 = 1
	testUser             = &userdomain.User{ID: testUserID}
	testDurations        = [][]time.Time{{time.Date(2025, 10, 10, 0, 0, 0, 0, time.UTC), time.Date(2025, 10, 11, 0, 0, 0, 0, time.UTC)}}
	testBooking          = &booking.Booking{UserID: testUserID, HotelID: testHotelName, RoomID: testRoomNumber, TotalPrice: 100, Currency: "USD", CheckIn: testDurations[0][0], CheckOut: testDurations[0][1], Status: booking.BookingStatusPaid, PaymentID: "1"}
	testProvider   *BookingProvider
)

type MockBookingRepo struct{}

func (m *MockBookingRepo) Save(ctx context.Context, booking *booking.Booking) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockBookingRepo) Delete(ctx context.Context, booking *booking.Booking) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockBookingRepo) IsIntersected(ctx context.Context, hotelName int64, hotelRoom int64, checkIn time.Time, checkOut time.Time) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockBookingRepo) GetByBookingInfo(ctx context.Context, bookingInfo *object.BookingInfo) (*booking.Booking, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockBookingRepo) GetByHotel(ctx context.Context, hotelID int64) ([]*booking.Booking, error) {
	if hotelID == testHotelName {
		return []*booking.Booking{testBooking}, nil
	}
	return make([]*booking.Booking, 0), nil
}

func (m *MockBookingRepo) GetByUser(ctx context.Context, userID string) ([]*booking.Booking, error) {
	if userID == testUserID {
		return []*booking.Booking{testBooking}, nil
	}
	return make([]*booking.Booking, 0), nil
}

func (m *MockBookingRepo) GetDurationsByRoom(ctx context.Context, hotelID int64, roomNumber int64) ([][]time.Time, error) {
	if hotelID == testHotelName && roomNumber == testRoomNumber {
		return testDurations, nil
	}
	return [][]time.Time{}, nil
}

func (m *MockBookingRepo) GetBookingsByStatus(ctx context.Context, status booking.BookingStatus) ([]*booking.Booking, error) {
	//TODO implement me
	panic("implement me")
}

type MockHotelRepo struct{}

func (m *MockHotelRepo) IsHotelier(ctx context.Context, userID string, hotelName int64) (bool, error) {
	if userID == testUserID && hotelName == testHotelName {
		return true, nil
	}

	return false, nil
}

func (m *MockHotelRepo) GetRoomInfo(ctx context.Context, hotelName int64, roomNumber int64) (*hotel.RoomInfo, error) {
	//TODO implement me
	panic("implement me")
}

func TestMain(m *testing.M) {
	testProvider = NewBookingProvider(&MockBookingRepo{}, &MockHotelRepo{})
	log.Println("Running tests...")
	code := m.Run()

	os.Exit(code)
}

func TestBookingProvider_GetOccupiedRoomDurations(t *testing.T) {
	ctx := context.Background()
	durations, err := testProvider.GetOccupiedRoomDurations(ctx, testHotelName, 2)
	assert.Assert(t, err == nil)
	assert.Equal(t, len(durations), 0, "Expected durations to be empty")
	durations, err = testProvider.GetOccupiedRoomDurations(ctx, testHotelName, testRoomNumber)
	assert.Assert(t, err == nil)
	assert.DeepEqual(t, testDurations, durations)
}

func TestBookingProvider_GetBookingsByHotelier(t *testing.T) {
	ctx := context.Background()
	_, err := testProvider.GetBookingsByHotelier(ctx, &userdomain.User{ID: "2"}, testHotelName)
	assert.Assert(t, err != nil)
	bookings, err := testProvider.GetBookingsByHotelier(ctx, testUser, testRoomNumber)
	assert.Assert(t, err == nil)
	assert.Assert(t, len(bookings) != 0)
	assert.DeepEqual(t, bookings[0].HotelID, testBooking.HotelID)
}

func TestBookingProvider_GetBookingsByUser(t *testing.T) {
	ctx := context.Background()
	bookings, err := testProvider.GetBookingsByUser(ctx, &userdomain.User{ID: "2"})
	assert.Assert(t, err == nil)
	assert.Assert(t, len(bookings) == 0)
	bookings, err = testProvider.GetBookingsByUser(ctx, testUser)
	assert.Assert(t, err == nil)
	assert.Assert(t, len(bookings) != 0)
}
