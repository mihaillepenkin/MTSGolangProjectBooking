package service

import (
	"database/sql"
	"hotel_service/internal/application/dto/request"
	"hotel_service/internal/application/dto/response"
)

type HotelService interface {
	Initialize(db *sql.DB)
	GetAllHotels() (response.AllHotelsInfoResponseDto, error)
	AddHotelInfo(hotelInfo *request.HotelInfoAdditionRequestDto) (response.HotelInfoResponseDto, error)
	UpdateHotelInfo(newHotelInfo *request.HotelInfoUpdateRequestDto) (response.HotelInfoResponseDto, error)
	GetHotelById(id int64) (response.HotelInfoResponseDto, error)
}
