package db

import (
	"database/sql"
	"errors"
	"hotel_service/internal/application/dto/request"
	"hotel_service/internal/domain/entity"
	"log/slog"
)

type HotelRepository struct {
	Db *sql.DB
}

func (hr *HotelRepository) Initialize(db *sql.DB) {
	hr.Db = db
}

func (hr *HotelRepository) CheckIfHotelExists(hotelName string, hotelLocation string) error {
	tx, err := hr.Db.Begin()
	if err != nil {
		return err
	}
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			slog.Error(err.Error())
		}
	}(tx)

	row := tx.QueryRow("SELECT name, location FROM hotels WHERE name = $1 AND location = $2", hotelName, hotelLocation)
	var name string
	var location string
	err = row.Scan(&name, &location)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (hr *HotelRepository) CheckHotelOwner(hotelId int64, userId int64) error {
	tx, err := hr.Db.Begin()
	if err != nil {
		return err
	}
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			slog.Error(err.Error())
		}
	}(tx)

	row := tx.QueryRow("SELECT owner_id FROM hotels WHERE id = $1", hotelId)
	var ownerId int64
	err = row.Scan(&ownerId)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	if ownerId != userId {
		return errors.New("user is not hotel owner")
	}

	return nil
}

func (hr *HotelRepository) GetAllHotels() ([]entity.Hotel, error) {
	tx, err := hr.Db.Begin()
	if err != nil {
		return nil, err
	}
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			slog.Error(err.Error())
		}
	}(tx)

	hotels := make([]entity.Hotel, 0)
	rows, err := tx.Query("SELECT id, name, description, location, owner_id FROM hotels")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var hotel entity.Hotel
		err := rows.Scan(&hotel.Id, &hotel.Name, &hotel.Description, &hotel.Location, &hotel.OwnerId)
		if err != nil {
			return nil, err
		}
		hotels = append(hotels, hotel)
	}

	rooms := make(map[int64][]entity.Room)
	rows, err = tx.Query("SELECT id, number, price, hotel_id FROM rooms")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var room entity.Room
		err := rows.Scan(&room.Id, &room.Number, &room.Price, &room.HotelId)
		if err != nil {
			return nil, err
		}
		rooms[room.HotelId] = append(rooms[room.HotelId], room)
	}

	for _, hotel := range hotels {
		hotel.Rooms = rooms[hotel.Id]
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return hotels, nil
}

func (hr *HotelRepository) AddHotelInfo(name string, description string, location string, ownerId int64, rooms []request.Room) (entity.Hotel, error) {
	tx, err := hr.Db.Begin()
	if err != nil {
		return entity.Hotel{}, err
	}
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			slog.Error(err.Error())
		}
	}(tx)

	hotel := entity.Hotel{Name: name, Description: description, Location: location, OwnerId: ownerId}
	result, err := tx.Exec("INSERT INTO hotels (name, description, location, owner_id) VALUES ($1, $2, $3, $4)", hotel.Name, hotel.Description, hotel.Location, hotel.OwnerId)
	if err != nil {
		return entity.Hotel{}, err
	}
	hotel.Id, err = result.LastInsertId()
	if err != nil {
		return entity.Hotel{}, err
	}

	roomsEnt := make([]entity.Room, 0)
	for _, room := range rooms {
		roomEnt := entity.Room{Number: room.Number, Price: room.Price, HotelId: hotel.Id}
		result, err := tx.Exec("INSERT INTO rooms (number, price, hotel_id) VALUES ($1, $2, $3)", roomEnt.Number, roomEnt.Price, roomEnt.HotelId)
		if err != nil {
			return entity.Hotel{}, err
		}
		roomEnt.Id, err = result.LastInsertId()
		if err != nil {
			return entity.Hotel{}, err
		}
		roomsEnt = append(roomsEnt, roomEnt)
	}
	hotel.Rooms = roomsEnt

	err = tx.Commit()
	if err != nil {
		return entity.Hotel{}, err
	}

	return hotel, nil
}

func (hr *HotelRepository) UpdateHotelInfo(id int64, newName string, newDescription string, newLocation string, newOwnerId int64, newRooms []request.RoomUpd) (entity.Hotel, error) {
	tx, err := hr.Db.Begin()
	if err != nil {
		return entity.Hotel{}, err
	}
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			slog.Error(err.Error())
		}
	}(tx)

	hotel := entity.Hotel{Id: id, Name: newName, Description: newDescription, Location: newLocation, OwnerId: newOwnerId}
	_, err = tx.Exec("UPDATE hotels SET name=$1, description=$2, location=$3, owner_id=$4 WHERE id=$5", hotel.Name, hotel.Description, hotel.Location, hotel.OwnerId, hotel.Id)
	if err != nil {
		return entity.Hotel{}, err
	}

	roomsEnt := make([]entity.Room, 0)
	for _, room := range newRooms {
		roomEnt := entity.Room{Id: room.Id, Number: room.Number, Price: room.Price, HotelId: hotel.Id}
		_, err := tx.Exec("UPDATE rooms SET number=$1, price=$2 WHERE id=$3", roomEnt.Number, roomEnt.Price, roomEnt.Id)
		if err != nil {
			return entity.Hotel{}, err
		}
		roomsEnt = append(roomsEnt, roomEnt)
	}
	hotel.Rooms = roomsEnt

	err = tx.Commit()
	if err != nil {
		return entity.Hotel{}, err
	}

	return hotel, nil
}
