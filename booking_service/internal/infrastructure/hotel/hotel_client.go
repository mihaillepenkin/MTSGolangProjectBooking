package hotel

import (
	"context"
	"log/slog"

	grpchotel "github.com/Vlad-Ali/MTSGolangProjectBooking-protos/gen/proto/hotel"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/hotel"
	error2 "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/hotel/error"
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

func (h *HotelClient) IsHotelier(ctx context.Context, userID string, hotelID int64) (bool, error) {
	request := &grpchotel.IsHotelierRequest{UserID: userID, HotelID: hotelID}

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

func (h *HotelClient) GetRoomInfo(ctx context.Context, hotelID int64, roomNumber int64) (*hotel.RoomInfo, error) {
	request := &grpchotel.RoomInfoRequest{HotelID: hotelID, RoomNumber: roomNumber}

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
