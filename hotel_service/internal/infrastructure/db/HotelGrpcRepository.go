package db

import (
	"database/sql"
	"log/slog"
	"strconv"
)

type HotelGrpcRepository struct {
	Db *sql.DB
}

func (hgr *HotelGrpcRepository) Initialize(db *sql.DB) {
	hgr.Db = db
}

func (hgr *HotelGrpcRepository) IsHotelier(hotelId int64, userId string) (bool, error) {
	tx, err := hgr.Db.Begin()
	if err != nil {
		return false, err
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
		return false, err
	}
	err = tx.Commit()
	if err != nil {
		return false, err
	}

	if strconv.FormatInt(ownerId, 10) != userId {
		return false, nil
	}

	return true, nil
}

func (hgr *HotelGrpcRepository) GetRoomInfo(hotelId int64, roomId int64) (int, int, string, string, error) {
	tx, err := hgr.Db.Begin()
	if err != nil {
		return 0, 0, "", "", err
	}
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			slog.Error(err.Error())
		}
	}(tx)

	row := tx.QueryRow("SELECT name FROM hotels WHERE id = $1", hotelId)
	var hotelName string
	err = row.Scan(&hotelName)
	if err != nil {
		return 0, 0, "", "", err
	}
	row = tx.QueryRow("SELECT number, price FROM rooms WHERE id = $1", roomId)
	var roomNumber int
	var price int
	err = row.Scan(&roomNumber, &price)
	if err != nil {
		return 0, 0, "", "", err
	}

	err = tx.Commit()
	if err != nil {
		return 0, 0, "", "", err
	}

	return roomNumber, price, "RUB", hotelName, nil
}
