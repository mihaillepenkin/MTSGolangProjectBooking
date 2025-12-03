package hotel

import (
	"context"
	"log/slog"

	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/hotel"
	error2 "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/hotel/error"
	grpchotel "github.com/mihaillepenkin/MTSGolangProjectBooking/protos/gen/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type HotelClient struct {
	client grpchotel.HotelClient
	logger *slog.Logger
}

func NewHotelClient(conn *grpc.ClientConn) *HotelClient {
	return &HotelClient{client: grpchotel.NewHotelClient(conn), logger: slog.Default().With("component", "hotel_client")}
}

func (h *HotelClient) IsHotelier(ctx context.Context, userID string, hotelName string) (bool, error) {
	request := &grpchotel.IsHotelierRequest{UserID: userID, HotelName: hotelName}

	response, err := h.client.IsHotelier(ctx, request)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			h.logger.Error("Failed to call client.IsHotelier()", "error", err)
			return false, err
		}

		switch st.Code() {
		case codes.InvalidArgument:
			h.logger.Error("Invalid argument", "error", st.Message())
			return false, error2.ErrHotelierInvalidArgument
		default:
			h.logger.Error("Unknown error", "error", st.Message())
			return false, err
		}
	}

	return response.IsHotelier, nil
}

func (h *HotelClient) GetRoomInfo(ctx context.Context, hotelName string, roomNumber string) (*hotel.RoomInfo, error) {
	request := &grpchotel.RoomInfoRequest{HotelName: hotelName, RoomNumber: roomNumber}

	response, err := h.client.GetRoomInfo(ctx, request)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			h.logger.Error("Failed to call client.GetRoomInfo()", "error", err)
			return nil, err
		}

		switch st.Code() {
		case codes.InvalidArgument:
			h.logger.Error("Invalid argument", "error", st.Message())
			return nil, error2.ErrHotelRoomInvalidArgument
		case codes.NotFound:
			h.logger.Error("Room is not found", "error", st.Message())
			return nil, error2.ErrHotelRoomIsNotFound
		default:
			h.logger.Error("Unknown error", "error", st.Message())
			return nil, err
		}
	}

	roomInfo := &hotel.RoomInfo{
		Currency: response.Currency,
		Amount:   response.Amount,
	}

	return roomInfo, nil
}
