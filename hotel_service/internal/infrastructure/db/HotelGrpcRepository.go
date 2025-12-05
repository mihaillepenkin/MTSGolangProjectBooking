package db

import (
	"database/sql"
	"errors"
	"log/slog"
	"strconv"
)

type HotelGrpcRepository struct {
	Db *sql.DB
}

func (hgr *HotelGrpcRepository) Initialize(db *sql.DB) {
	hgr.Db = db
}

func (hgr *HotelGrpcRepository) IsHotelier(hotelName string, userId string) (bool, error) {
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

	row := tx.QueryRow("SELECT owner_id FROM hotels WHERE name = $1", hotelName)
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
		return false, errors.New("user is not hotel owner")
	}

	return true, nil
}

func (hgr *HotelGrpcRepository) GetRoomInfo(hotelName string, roomNumber string) (int64, string, error) {
	tx, err := hgr.Db.Begin()
	if err != nil {
		return 0, "", err
	}
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			slog.Error(err.Error())
		}
	}(tx)

	row := tx.QueryRow("SELECT id FROM hotels WHERE name = $1", hotelName)
	var id int64
	err = row.Scan(&id)
	if err != nil {
		return 0, "", err
	}
	row = tx.QueryRow("SELECT price FROM rooms WHERE hotel_id = $1", id)
	var price int64
	err = row.Scan(&price)
	if err != nil {
		return 0, "", err
	}

	err = tx.Commit()
	if err != nil {
		return 0, "", err
	}

	return price, "RUB", nil
}
