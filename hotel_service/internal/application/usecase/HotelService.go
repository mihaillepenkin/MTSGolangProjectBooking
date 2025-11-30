package usecase

import (
	"database/sql"
	"errors"
	"hotel_service/internal/application/dto/request"
	"hotel_service/internal/application/dto/response"
	"hotel_service/internal/infrastructure/db"
	"log/slog"
)

type HotelService struct {
	hotelRepository db.HotelRepository
}

func (hs *HotelService) Initialize(db *sql.DB) {
	hs.hotelRepository.Initialize(db)
}

func (hs *HotelService) GetAllHotels() (response.AllHotelsInfoResponseDto, error) {
	hotels, err := hs.hotelRepository.GetAllHotels()
	if err != nil {
		slog.Error("Ошибка в HotelRepository, метод GetAllHotels: " + err.Error())
		return response.AllHotelsInfoResponseDto{Message: "Ошибка при получении информации о всех отелях", Error: "500"}, errors.New("ошибка в ходе работы с БД в репозитории")
	}

	result := make([]response.Hotel, len(hotels))
	for i, hotel := range hotels {
		result[i] = response.Hotel{Id: hotel.Id, Name: hotel.Name, Location: hotel.Location, OwnerId: hotel.OwnerId}
		rooms := make([]response.Room, len(hotel.Rooms))
		for j, room := range hotel.Rooms {
			rooms[j] = response.Room{Id: room.Id, Number: room.Number, Price: room.Price}
		}
		result[i].Rooms = rooms
	}

	return response.AllHotelsInfoResponseDto{Hotels: result, Message: "", Error: ""}, nil
}

func (hs *HotelService) AddHotelInfo(hotelInfo *request.HotelInfoAdditionRequestDto) (response.HotelInfoResponseDto, error) {
	hotel, err := hs.hotelRepository.AddHotelInfo(hotelInfo.Name, hotelInfo.Location, hotelInfo.OwnerId, hotelInfo.Rooms)
	if err != nil {
		slog.Error("Ошибка в HotelRepository, метод AddHotelInfo: " + err.Error())
		return response.HotelInfoResponseDto{Message: "Ошибка при добавлении информации об отеле", Error: "500"}, errors.New("ошибка в ходе работы с БД в репозитории")
	}

	rooms := make([]response.Room, len(hotel.Rooms))
	for j, room := range hotel.Rooms {
		rooms[j] = response.Room{Id: room.Id, Number: room.Number, Price: room.Price}
	}

	return response.HotelInfoResponseDto{Id: hotel.Id, Name: hotel.Name, Location: hotel.Location, OwnerId: hotel.OwnerId, Rooms: rooms, Message: "", Error: ""}, nil
}

func (hs *HotelService) UpdateHotelInfo(newHotelInfo *request.HotelInfoUpdateRequestDto) (response.HotelInfoResponseDto, error) {
	hotelUpd, err := hs.hotelRepository.UpdateHotelInfo(newHotelInfo.Id, newHotelInfo.NewName, newHotelInfo.NewLocation, newHotelInfo.NewOwnerId, newHotelInfo.NewRooms)
	if err != nil {
		slog.Error("Ошибка в HotelRepository, метод UpdateHotelInfo: " + err.Error())
		return response.HotelInfoResponseDto{Message: "Ошибка при обновлении информации об отеле", Error: "500"}, errors.New("ошибка в ходе работы с БД в репозитории")
	}

	roomsUpd := make([]response.Room, len(hotelUpd.Rooms))
	for i, room := range hotelUpd.Rooms {
		roomsUpd[i] = response.Room{Id: room.Id, Number: room.Number, Price: room.Price}
	}

	return response.HotelInfoResponseDto{Id: hotelUpd.Id, Name: hotelUpd.Name, Location: hotelUpd.Location, OwnerId: hotelUpd.OwnerId, Rooms: roomsUpd, Message: "", Error: ""}, nil
}

func (hs *HotelService) GetHotelById(id int64) (response.HotelInfoResponseDto, error) {
	return response.HotelInfoResponseDto{}, nil
}
