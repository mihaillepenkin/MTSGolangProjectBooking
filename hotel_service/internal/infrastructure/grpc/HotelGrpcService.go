package grpc2

import (
	"context"
	"database/sql"
	"errors"
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

func (hgs *HotelGrpcService) Initialize(database *sql.DB) {
	hgs.hotelGrpcRepository = &db.HotelGrpcRepository{}
	hgs.hotelGrpcRepository.Initialize(database)
}

func (hgs *HotelGrpcService) IsHotelier(ctx context.Context, r *hotel.IsHotelierRequest) (*hotel.IsHotelierResponse, error) {
	isHotelier, err := hgs.hotelGrpcRepository.IsHotelier(r.HotelID, r.UserID)
	if err != nil {
		slog.Error("Ошибка в HotelGrpcRepository, метод IsHotelier: " + err.Error())
		return nil, status.Error(codes.Internal, "Ошибка при проверке личности владельца отеля")
	}

	return &hotel.IsHotelierResponse{IsHotelier: isHotelier}, nil
}

func (hgs *HotelGrpcService) GetRoomInfo(ctx context.Context, r *hotel.RoomInfoRequest) (*hotel.RoomInfoResponse, error) {
	roomNumber, price, currency, hotelName, err := hgs.hotelGrpcRepository.GetRoomInfo(r.HotelID, r.RoomID)

	if errors.Is(err, sql.ErrNoRows) {
		slog.Error("Не найден информация о номере", "hotelID", r.HotelID, "roomID", r.RoomID)
		return nil, status.Error(codes.NotFound, "Room information not found")
	}

	if err != nil {
		slog.Error("Ошибка в HotelGrpcRepository, метод GetRoomInfo: " + err.Error())
		return nil, status.Error(codes.Internal, "Ошибка при получении информации о номере в отеле")
	}

	return &hotel.RoomInfoResponse{Amount: int32(price), Currency: currency, HotelName: hotelName, RoomNumber: int32(roomNumber)}, nil
}
