package tests

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"hotel_service/internal/application/dto/request"
	"hotel_service/internal/domain/entity"
	"hotel_service/internal/infrastructure/db"
	"log"
	"log/slog"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"gotest.tools/v3/assert"
)

var (
	testDb        *sql.DB
	testRepo      *db.HotelRepository
	testHotelInfo = entity.Hotel{
		Id:          1,
		Name:        "Hotel 1",
		Description: "Hotel Description",
		Location:    "Moscow",
		OwnerId:     13,
		Rooms:       []entity.Room{{Id: 1, Number: 567, Price: 789, HotelId: 1}},
	}
	testHotelInfoUpd = entity.Hotel{
		Id:          1,
		Name:        "Hotel 1",
		Description: "Hotel Description",
		Location:    "Moscow, Main street",
		OwnerId:     13,
		Rooms:       []entity.Room{{Id: 1, Number: 567, Price: 897, HotelId: 1}},
	}
	testHotelList = []entity.Hotel{testHotelInfoUpd}
)

func setupTests() error {
	container, err := postgres.Run(
		context.Background(),
		"postgres:15-alpine",
		postgres.WithDatabase("testDb"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		postgres.WithInitScripts("../migrations/001__create_tables.up.sql"),
	)
	if err != nil {
		slog.Error(err.Error())
		return errors.New("error running postgres container")
	}

	connStr, err := container.ConnectionString(context.Background())
	if err != nil {
		slog.Error(err.Error())
		return errors.New("error getting database connection string")
	}
	connStr += "sslmode=disable"
	slog.Info(connStr)
	dbT, err := sql.Open("postgres", connStr)
	if err != nil {
		slog.Error(err.Error())
		return errors.New("error connecting to postgres")
	}

	testDb = dbT
	testRepo = &db.HotelRepository{Db: testDb}

	return nil
}

func TestMain(m *testing.M) {
	err := setupTests()
	if err != nil {
		slog.Error(err.Error())
		log.Fatal("error setting up tests")
	}

	slog.Info("running tests...")
	err = testDb.Ping()
	if err != nil {
		slog.Error(err.Error())
	}
	code := m.Run()

	err = testDb.Close()
	if err != nil {
		slog.Error("error closing database: err", "err", err.Error())
	}

	os.Exit(code)
}

func TestHotelRepositoryAddHotelInfo(t *testing.T) {
	hotelInfo, err := testRepo.AddHotelInfo("Hotel 1", "Hotel Description", "Moscow", 13, []request.Room{{Number: 567, Price: 789}})
	if err != nil {
		t.Fatalf("error adding hotelInfo: %v", err)
	}

	t.Log(hotelInfo, testHotelInfo)

	assert.DeepEqual(t, testHotelInfo, hotelInfo)
}

func TestHotelRepositoryUpdateHotelInfo(t *testing.T) {
	hotelInfoUpd, err := testRepo.UpdateHotelInfo(1, "Hotel 1", "Hotel Description", "Moscow, Main street", 13, []request.RoomUpd{{Id: 1, Number: 567, Price: 897}})
	if err != nil {
		t.Fatalf("error updating hotelInfo: %v", err)
	}

	t.Log(hotelInfoUpd, testHotelInfoUpd)

	assert.DeepEqual(t, testHotelInfoUpd, hotelInfoUpd)
}

func TestHotelRepositoryGetAllHotels(t *testing.T) {
	hotels, err := testRepo.GetAllHotels()
	if err != nil {
		t.Fatalf("error getting hotels: %v", err)
	}

	t.Log(hotels, testHotelList)

	assert.DeepEqual(t, testHotelList, hotels)
}

func TestHotelRepositoryCheckHotelOwner(t *testing.T) {
	isHotelOwner, err := testRepo.CheckHotelOwner(1, 17)
	if err != nil {
		t.Fatalf("error checking hotel owner: %v", err)
	}

	t.Log(isHotelOwner, false)

	assert.Equal(t, false, isHotelOwner, fmt.Sprintf("isHotelOwner must be %v", false))
}

func TestHotelRepositoryCheckIfHotelExists(t *testing.T) {
	ifHotelExists, err := testRepo.CheckIfHotelExists("Hotel 1", "Moscow, Main street")
	if err != nil {
		t.Fatalf("error checking if hotel exists: %v", err)
	}

	t.Log(ifHotelExists, true)

	assert.Equal(t, true, ifHotelExists, fmt.Sprintf("ifHotelExists must be %v", true))
}
