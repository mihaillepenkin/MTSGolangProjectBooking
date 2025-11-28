package hotel

import "context"

type Repository interface {
	IsHotelier(ctx context.Context, userID string, hotelName string) (bool, error)
	GetRoomInfo(ctx context.Context, hotelName string, roomNumber string) (*RoomInfo, error)
}
