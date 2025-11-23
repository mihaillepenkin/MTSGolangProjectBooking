package usecase

import (
	"database/sql"
	"hotel_service/internal/application/dto/request"
	"hotel_service/internal/application/dto/response"
	"hotel_service/internal/infrastructure/db"
)

type HotelService struct {
	hotelRepository db.HotelRepository
}

func (hs *HotelService) Initialize(db *sql.DB) {
	hs.hotelRepository.Initialize(db)
}

func (hs *HotelService) GetAllHotels() (response.AllHotelsInfoResponseDto, error) {
	return response.AllHotelsInfoResponseDto{}, nil
}

func (hs *HotelService) AddHotelInfo(hotelInfo *request.HotelInfoAdditionRequestDto) (response.HotelInfoResponseDto, error) {
	return response.HotelInfoResponseDto{}, nil
}

func (hs *HotelService) UpdateHotelInfo(newHotelInfo *request.HotelInfoUpdateRequestDto) (response.HotelInfoResponseDto, error) {
	return response.HotelInfoResponseDto{}, nil
}

func (hs *HotelService) GetHotelById(id int64) (response.HotelInfoResponseDto, error) {
	return response.HotelInfoResponseDto{}, nil
}
