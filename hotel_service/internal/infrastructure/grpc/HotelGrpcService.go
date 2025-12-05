package grpc2

import (
	"context"
	"database/sql"
	"hotel_service/internal/infrastructure/db"
	"log/slog"

	"github.com/Vlad-Ali/MTSGolangProjectBooking-protos/gen/proto/hotel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type HotelGrpcService struct {
	hotelGrpcRepository *db.HotelGrpcRepository
	hotel.UnimplementedHotelServer
}

func (hgs *HotelGrpcService) Initialize(db *sql.DB) {
	hgs.hotelGrpcRepository.Initialize(db)
}

func (hgs *HotelGrpcService) IsHotelier(ctx context.Context, r *hotel.IsHotelierRequest) (*hotel.IsHotelierResponse, error) {
	isHotelier, err := hgs.hotelGrpcRepository.IsHotelier(r.HotelName, r.UserID)
	if err != nil {
		if err.Error() == "user is not hotel owner" {
			slog.Error("Ошибка в HotelGrpcService, метод IsHotelier: " + err.Error())
			return nil, status.Error(codes.PermissionDenied, "Ошибка при проверке личности владельца отеля")
		}
		slog.Error("Ошибка в HotelGrpcRepository, метод IsHotelier: " + err.Error())
		return nil, status.Error(codes.Internal, "Ошибка при проверке личности владельца отеля")
	}

	return &hotel.IsHotelierResponse{IsHotelier: isHotelier}, nil
}

func (hgs *HotelGrpcService) GetRoomInfo(ctx context.Context, r *hotel.RoomInfoRequest) (*hotel.RoomInfoResponse, error) {
	amount, currency, err := hgs.hotelGrpcRepository.GetRoomInfo(r.HotelName, r.RoomNumber)
	if err != nil {
		slog.Error("Ошибка в HotelGrpcRepository, метод GetRoomInfo: " + err.Error())
		return nil, status.Error(codes.Internal, "Ошибка при получении информации о номере в отеле")
	}

	return &hotel.RoomInfoResponse{Amount: amount, Currency: currency}, nil
}
