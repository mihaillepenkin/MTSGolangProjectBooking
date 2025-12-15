package booking

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/adapter/booking/request"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/adapter/userkey"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking"
	error2 "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking/error"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking/object"
	error3 "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/hotel/error"
	userdomain "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/user"
)

var (
	testHotelID    int64 = 1
	testRoomNumber int64 = 1
	testDurations        = []time.Time{time.Date(2025, 10, 10, 0, 0, 0, 0, time.UTC), time.Date(2025, 10, 12, 0, 0, 0, 0, time.UTC)}
	testUser             = &userdomain.User{ID: uuid.New().String(), Role: userdomain.RoleClient}
	testServer     *httptest.Server
	testClient     *http.Client
)

const (
	BookingPrefix            = "/booking"
	BookRoomEndpoint         = BookingPrefix
	GetHotelBookingsEndpoint = BookingPrefix + "/bookings" + "/hotelier"
	GetUserBookingsEndpoint  = BookingPrefix + "/client"
	GetRoomDurationsEndpoint = BookingPrefix + "/occupied"
)

type MockBookingSaver struct{}

func (m *MockBookingSaver) BookRoom(ctx context.Context, bookingInfo *object.BookingInfo) (string, error) {
	if bookingInfo.CheckOut.Before(bookingInfo.CheckIn) {
		return "", error2.ErrTimeDurationIsNotValid
	}

	if bookingInfo.CheckIn != testDurations[0] || bookingInfo.CheckOut != testDurations[1] {
		return "", error2.ErrBookingIsIntersected
	}

	if bookingInfo.HotelID != testHotelID || bookingInfo.RoomID != testRoomNumber {
		return "", error3.ErrHotelRoomIsNotFound
	}

	return "url", nil
}

func (m *MockBookingSaver) DeleteBooking(ctx context.Context, bookingInfo *object.BookingInfo) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockBookingSaver) ConfirmBooking(ctx context.Context, bookingInfo *object.BookingInfo) error {
	//TODO implement me
	panic("implement me")
}

type MockBookingProvider struct{}

func (m *MockBookingProvider) GetBookingsByUser(ctx context.Context, user *userdomain.User) ([]*booking.Booking, error) {
	return make([]*booking.Booking, 0), nil
}

func (m *MockBookingProvider) GetBookingsByHotelier(ctx context.Context, user *userdomain.User, hotelID int64) ([]*booking.Booking, error) {
	if user.ID != testUser.ID {
		return nil, error2.ErrHotelierIsNotValid
	}

	return make([]*booking.Booking, 0), nil
}

func (m *MockBookingProvider) GetOccupiedRoomDurations(ctx context.Context, hotelID int64, roomNumber int64) ([][]time.Time, error) {
	return [][]time.Time{testDurations}, nil
}

func testMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), userkey.UserKey{}, testUser)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func TestMain(m *testing.M) {
	bookHandler := NewBookingHandler(&MockBookingSaver{}, &MockBookingProvider{})

	mux := http.NewServeMux()
	mux.HandleFunc(BookRoomEndpoint, bookHandler.BookRoom)
	mux.HandleFunc(GetHotelBookingsEndpoint, bookHandler.GetHotelBookings)
	mux.HandleFunc(GetUserBookingsEndpoint, bookHandler.GetUserBookings)
	mux.HandleFunc(GetRoomDurationsEndpoint, bookHandler.GetOccupiedRoomDurations)

	testServer = httptest.NewServer(testMiddleware(mux))
	testClient = testServer.Client()

	log.Println("Running test...")
	code := m.Run()

	os.Exit(code)
}

func TestBookingHandler_BookRoom(t *testing.T) {
	tests := []struct {
		name           string
		request        request.BookRoomRequest
		expectedStatus int
	}{
		{
			name:           "Invalid duration",
			request:        request.BookRoomRequest{HotelID: testHotelID, RoomID: testRoomNumber, CheckIn: time.Date(2025, 10, 11, 0, 0, 0, 0, time.UTC), CheckOut: time.Date(2025, 10, 10, 0, 0, 0, 0, time.UTC)},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid intersected duration",
			request:        request.BookRoomRequest{HotelID: testHotelID, RoomID: testRoomNumber, CheckIn: time.Date(2025, 10, 11, 0, 0, 0, 0, time.UTC), CheckOut: time.Date(2025, 10, 12, 0, 0, 0, 0, time.UTC)},
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "Hotel room is not found",
			request:        request.BookRoomRequest{HotelID: testHotelID, RoomID: 2, CheckIn: testDurations[0], CheckOut: testDurations[1]},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Hotel room is found",
			request:        request.BookRoomRequest{HotelID: testHotelID, RoomID: testRoomNumber, CheckIn: testDurations[0], CheckOut: testDurations[1]},
			expectedStatus: http.StatusOK,
		},
	}

	for _, test := range tests {
		jsonData, err := json.Marshal(&test.request)
		if err != nil {
			t.Fatal("Error marshalling json", err)
		}

		req, err := http.NewRequest(http.MethodPost, testServer.URL+BookRoomEndpoint, bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatal("Error creating request", err)
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := testClient.Do(req)
		if err != nil {
			t.Fatal("Error sending request", err)
		}
		if resp.StatusCode != test.expectedStatus {
			t.Errorf("Invalid status code. Expected %d, got %d", test.expectedStatus, resp.StatusCode)
		}
	}
}

func TestBookingHandler_GetHotelBookings(t *testing.T) {
	tests := []struct {
		name           string
		hotelID        string
		expectedStatus int
	}{
		{
			name:           "Not authorized",
			hotelID:        "1",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, test := range tests {
		url := testServer.URL + GetHotelBookingsEndpoint + "?hotelID=" + test.hotelID
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			t.Fatal("Error creating request", err)
		}

		resp, err := testClient.Do(req)
		if err != nil {
			t.Fatal("Error sending request", err)
		}

		if resp.StatusCode != test.expectedStatus {
			t.Errorf("Invalid status code. Expected %d, got %d", test.expectedStatus, resp.StatusCode)
		}
	}
}

func TestBookingHandler_GetUserBookings(t *testing.T) {
	tests := []struct {
		name           string
		expectedStatus int
	}{
		{name: "Success", expectedStatus: http.StatusOK},
	}
	for _, test := range tests {
		url := testServer.URL + GetUserBookingsEndpoint
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			t.Fatal("Error creating request", err)
		}

		resp, err := testClient.Do(req)
		if err != nil {
			t.Fatal("Error sending request", err)
		}

		if resp.StatusCode != test.expectedStatus {
			t.Errorf("Invalid status code. Expected %d, got %d", test.expectedStatus, resp.StatusCode)
		}
	}
}

func TestBookingHandler_GetRoomDurations(t *testing.T) {
	tests := []struct {
		name           string
		hotelID        string
		roomNumber     string
		expectedStatus int
	}{
		{name: "Invalid hotelID", expectedStatus: http.StatusBadRequest},
		{name: "Success", hotelID: "1", roomNumber: "1", expectedStatus: http.StatusOK},
	}

	for _, test := range tests {
		url := testServer.URL + GetRoomDurationsEndpoint + "?hotelID=" + test.hotelID + "&roomID=" + test.roomNumber
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			t.Fatal("Error creating request", err)
		}

		resp, err := testClient.Do(req)
		if err != nil {
			t.Fatal("Error sending request", err)
		}

		if resp.StatusCode != test.expectedStatus {
			t.Errorf("Invalid status code. Expected %d, got %d", test.expectedStatus, resp.StatusCode)
		}
	}
}
