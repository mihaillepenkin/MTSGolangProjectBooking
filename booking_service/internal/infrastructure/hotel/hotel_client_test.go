package hotel

import (
	"context"
	"log"
	"log/slog"
	"os"
	"testing"

	"github.com/Vlad-Ali/MTSGolangProjectBooking-protos/gen/proto/hotel"
	"github.com/google/uuid"
	error2 "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/hotel/error"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gotest.tools/v3/assert"
)

var (
	testHotelClient *HotelClient
	userID                = uuid.New().String()
	hotelID         int64 = 1
	roomNumber      int64 = 1
)

type MockClient struct{}

func (*MockClient) IsHotelier(ctx context.Context, in *hotel.IsHotelierRequest, opts ...grpc.CallOption) (*hotel.IsHotelierResponse, error) {
	if in.UserID == "" || in.HotelID < 0 {
		return nil, status.Error(codes.InvalidArgument, "UserID is invalid")
	}

	if in.HotelID == hotelID && userID == in.UserID {
		return &hotel.IsHotelierResponse{IsHotelier: true}, nil
	}

	return &hotel.IsHotelierResponse{IsHotelier: false}, nil
}

func (*MockClient) GetRoomInfo(ctx context.Context, in *hotel.RoomInfoRequest, opts ...grpc.CallOption) (*hotel.RoomInfoResponse, error) {
	if in.RoomNumber < 0 || in.HotelID < 0 {
		return nil, status.Error(codes.InvalidArgument, "Request is invalid")
	}

	if in.HotelID == hotelID && in.RoomNumber == roomNumber {
		return &hotel.RoomInfoResponse{Currency: "USD", Amount: 100}, nil
	}

	return nil, status.Error(codes.NotFound, "Room not found")
}

func TestMain(m *testing.M) {
	testHotelClient = &HotelClient{client: &MockClient{}, logger: slog.Default().With("component", "hotel_client")}
	log.Println("Running tests...")
	code := m.Run()
	os.Exit(code)
}

func TestHotelClient_ShouldReturnInvalidArgumentErrForIsHotelier(t *testing.T) {
	ctx := context.Background()
	_, err := testHotelClient.IsHotelier(ctx, "", -1)
	assert.Equal(t, error2.ErrHotelierInvalidArgument, err, "expected InvalidArgumentErr")
}

func TestHotelClient_ShouldReturnFalseForIsHotelier(t *testing.T) {
	ctx := context.Background()
	isHotelier, err := testHotelClient.IsHotelier(ctx, userID, 2)
	assert.Assert(t, err == nil, "expected err to be nil")
	assert.Assert(t, isHotelier == false, "expected isHotelier to be false")
}

func TestHotelClient_ShouldReturnTrueIsHotelier(t *testing.T) {
	ctx := context.Background()
	isHotelier, err := testHotelClient.IsHotelier(ctx, userID, hotelID)
	assert.Assert(t, err == nil, "expected err to be nil")
	assert.Equal(t, isHotelier, true, "expected isHotelier to be true")
}

func TestHotelClient_ShouldReturnInvalidArgumentErrForGetRoomInfo(t *testing.T) {
	ctx := context.Background()
	_, err := testHotelClient.GetRoomInfo(ctx, -1, -1)
	assert.Equal(t, error2.ErrHotelRoomInvalidArgument, err, "expected InvalidArgumentErr")
}

func TestHotelClient_ShouldReturnNotFoundErrForGetRoomInfo(t *testing.T) {
	ctx := context.Background()
	_, err := testHotelClient.GetRoomInfo(ctx, hotelID, 2)
	assert.Equal(t, error2.ErrHotelRoomIsNotFound, err, "expected NotFoundErr")
}

func TestHotelClient_GetRoomInfo(t *testing.T) {
	ctx := context.Background()
	roomInfo, err := testHotelClient.GetRoomInfo(ctx, hotelID, roomNumber)
	assert.Assert(t, err == nil, "expected err to be nil")
	assert.Assert(t, roomInfo != nil, "expected roomInfo to not be nil")
}
