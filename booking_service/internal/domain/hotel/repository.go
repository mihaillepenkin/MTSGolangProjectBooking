package hotel

import "context"

type Repository interface {
	IsHotelier(ctx context.Context, userID string, hotelID int64) (bool, error)
	GetRoomInfo(ctx context.Context, hotelID int64, roomID int64) (*RoomInfo, error)
}
