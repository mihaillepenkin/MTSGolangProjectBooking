package repository

import (
	"database/sql"
	"hotel_service/internal/application/dto/request"
	"hotel_service/internal/domain/entity"
)

type HotelRepository interface {
	Initialize(db *sql.DB)
	GetAllHotels() ([]entity.Hotel, error)
	AddHotelInfo(name string, location string, ownerId int64, rooms []request.Room) (entity.Hotel, error)
	UpdateHotelInfo(id int64, newName string, newLocation string, newOwnerId int64, newRooms []request.RoomUpd) (entity.Hotel, error)
	GetHotelById(id int64) (entity.Hotel, error)
}
