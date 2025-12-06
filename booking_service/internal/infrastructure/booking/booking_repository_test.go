package booking

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	bookingdomain "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking"
	error2 "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking/error"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking/object"
	userdomain "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/user"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/usecase/case/transactionmanager"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"gotest.tools/v3/assert"
)

var (
	testDB      *sql.DB
	testRepo    *BookingRepository
	testBooking = bookingdomain.Booking{UserID: uuid.New().String(), HotelID: 1, RoomNumber: 1, TotalPrice: 100, Currency: "USD", CheckIn: time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC),
		CheckOut: time.Date(2025, 10, 2, 0, 0, 0, 0, time.UTC), Status: bookingdomain.BookingStatusUnpaid, PaymentID: uuid.New().String()}
)

func equalBookings(firstBooking bookingdomain.Booking, secondBooking bookingdomain.Booking) bool {
	firstBooking.CheckIn = firstBooking.CheckIn.UTC()
	firstBooking.CheckOut = firstBooking.CheckOut.UTC()
	secondBooking.CheckIn = secondBooking.CheckIn.UTC()
	secondBooking.CheckOut = secondBooking.CheckOut.UTC()
	return firstBooking == secondBooking
}

func setup() error {
	container, err := postgres.Run(context.Background(), "postgres:15-alpine",
		postgres.WithDatabase("testDB"), postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		postgres.WithOrderedInitScripts("../../../migrations/001_create_bookings_table.up.sql", "../../../migrations/002_add_no_overlapping_constraint.up.sql"))
	if err != nil {
		return fmt.Errorf("error running postgres container: %w", err)
	}

	time.Sleep(2 * time.Second)
	connStr, _ := container.ConnectionString(context.Background())
	connStr += "sslmode=disable"
	log.Println("ConnStr: ", connStr)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("error connecting to postgres: %w", err)
	}

	testDB = db
	testRepo = NewBookingRepository(testDB)
	return nil
}

func TestMain(m *testing.M) {
	if err := setup(); err != nil {
		log.Fatalf("error setting up tests: %v", err)
	}

	log.Println("Running tests...")
	code := m.Run()

	if err := testDB.Close(); err != nil {
		log.Printf("error closing db: %v", err)
	}

	os.Exit(code)
}

func TestBookingRepository_Save(t *testing.T) {
	booking := testBooking

	ctx := context.Background()
	tx, err := testDB.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("error starting transaction: %v", err)
	}

	defer tx.Rollback()

	txCtx := context.WithValue(ctx, transactionmanager.TxKey{}, tx)
	assert.Equal(t, nil, testRepo.Save(txCtx, &booking), fmt.Sprintf("testRepo.Save must be %v", nil))

	var userID string
	resErr := tx.QueryRowContext(txCtx, `
SELECT user_id FROM bookings WHERE user_id = $1`, booking.UserID).Scan(&userID)
	if resErr != nil {
		t.Fatalf("error querying bookings: %v", resErr)
	}

	assert.Assert(t, userID == booking.UserID, fmt.Sprintf("userID must be %v", booking.UserID))
}

func TestBookingRepository_UpdateBooking(t *testing.T) {
	booking := testBooking

	ctx := context.Background()
	tx, err := testDB.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("error starting transaction: %v", err)
	}

	defer tx.Rollback()
	txCtx := context.WithValue(ctx, transactionmanager.TxKey{}, tx)
	assert.Equal(t, nil, testRepo.Save(txCtx, &booking), fmt.Sprintf("testRepo.Save must be %v", nil))

	booking.Status = bookingdomain.BookingStatusPaid
	assert.Equal(t, nil, testRepo.Save(txCtx, &booking), fmt.Sprintf("testRepo.Save must be %v", nil))
	var status bookingdomain.BookingStatus
	resErr := tx.QueryRowContext(txCtx, `SELECT status FROM bookings WHERE user_id = $1`, booking.UserID).Scan(&status)
	if resErr != nil {
		t.Fatalf("error querying bookings: %v", resErr)
	}

	assert.Equal(t, booking.Status, status)

}

func TestBookingRepository_Delete(t *testing.T) {
	booking := testBooking
	ctx := context.Background()
	tx, err := testDB.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("error starting transaction: %v", err)
	}

	defer tx.Rollback()
	txCtx := context.WithValue(ctx, transactionmanager.TxKey{}, tx)
	assert.Equal(t, nil, testRepo.Save(txCtx, &booking), fmt.Sprintf("testRepo.Save must be %v", nil))
	assert.Equal(t, nil, testRepo.Delete(txCtx, &booking), fmt.Sprintf("testRepo.Delete must be %v", nil))

	var userID string
	resErr := tx.QueryRowContext(txCtx, `SELECT * FROM bookings`).Scan(&userID)
	assert.Assert(t, errors.Is(resErr, sql.ErrNoRows), "booking should be deleted")
}

func TestBookingRepository_ShouldReturnNotFoundErrForDelete(t *testing.T) {
	booking := testBooking
	ctx := context.Background()
	tx, err := testDB.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("error starting transaction: %v", err)
	}
	defer tx.Rollback()
	txCtx := context.WithValue(ctx, transactionmanager.TxKey{}, tx)
	assert.Equal(t, error2.ErrBookingIsNotFound, testRepo.Delete(txCtx, &booking), fmt.Sprintf("testRepo.Delete must be %v", error2.ErrBookingIsNotFound))
}

func TestBookingRepository_GetDurationsByRoom(t *testing.T) {
	booking := testBooking
	booking.Status = bookingdomain.BookingStatusPaid
	ctx := context.Background()
	tx, err := testDB.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("error starting transaction: %v", err)
	}

	defer tx.Rollback()

	txCtx := context.WithValue(ctx, transactionmanager.TxKey{}, tx)
	assert.Equal(t, nil, testRepo.Save(txCtx, &booking), fmt.Sprintf("testRepo.Save must be %v", nil))
	durations, err := testRepo.GetDurationsByRoom(txCtx, booking.HotelID, booking.RoomNumber)
	if err != nil {
		t.Fatalf("error getting durations: %v", err)
	}

	t.Log(durations, [][]time.Time{{booking.CheckIn, booking.CheckOut}})

	assert.DeepEqual(t, [][]time.Time{{booking.CheckIn, booking.CheckOut}}, durations)
}

func TestBookingRepository_IsIntersected(t *testing.T) {
	firstBooking := testBooking
	firstBooking.Status = bookingdomain.BookingStatusPaid

	secondBooking := testBooking
	secondBooking.CheckIn = firstBooking.CheckOut
	secondBooking.CheckOut = secondBooking.CheckIn.AddDate(0, 0, 1)
	thirdBooking := testBooking
	thirdBooking.CheckOut = thirdBooking.CheckIn.AddDate(0, 0, 5)
	ctx := context.Background()
	tx, err := testDB.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("error starting transaction: %v", err)
	}

	defer tx.Rollback()
	txCtx := context.WithValue(ctx, transactionmanager.TxKey{}, tx)
	assert.Equal(t, nil, testRepo.Save(txCtx, &firstBooking), fmt.Sprintf("testRepo.Save must be %v", nil))
	isIntersected, err := testRepo.IsIntersected(txCtx, secondBooking.HotelID, secondBooking.RoomNumber, secondBooking.CheckIn, secondBooking.CheckOut)
	if err != nil {
		t.Fatalf("error getting isIntersected: %v", err)
	}

	assert.Equal(t, false, isIntersected, "isIntersected must be false")

	isIntersected, err = testRepo.IsIntersected(txCtx, thirdBooking.HotelID, thirdBooking.RoomNumber, thirdBooking.CheckIn, thirdBooking.CheckOut)
	if err != nil {
		t.Fatalf("error getting isIntersected: %v", err)
	}

	assert.Equal(t, true, isIntersected, "isIntersected must be true")
}

func TestBookingRepository_GetByBookingInfo(t *testing.T) {
	booking := testBooking

	ctx := context.Background()
	tx, err := testDB.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("error starting transaction: %v", err)
	}

	defer tx.Rollback()
	txCtx := context.WithValue(ctx, transactionmanager.TxKey{}, tx)
	assert.Equal(t, nil, testRepo.Save(txCtx, &booking), fmt.Sprintf("testRepo.Save must be %v", nil))

	newBooking, err := testRepo.GetByBookingInfo(txCtx, &object.BookingInfo{HotelID: booking.HotelID, RoomNumber: booking.RoomNumber, CheckIn: booking.CheckIn, CheckOut: booking.CheckOut, User: userdomain.User{ID: booking.UserID}})
	if err != nil {
		t.Fatalf("error getting booking info: %v", err)
	}

	if newBooking == nil {
		t.Fatalf("booking info must not be nil")
	}

	booking.ID = newBooking.ID

	assert.Assert(t, equalBookings(booking, *newBooking), "booking must be equal to received booking")
}

func TestBookingRepository_ShouldReturnNotFoundErrForGetBookingInfo(t *testing.T) {
	booking := testBooking

	ctx := context.Background()
	tx, err := testDB.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("error starting transaction: %v", err)
	}

	defer tx.Rollback()

	txCtx := context.WithValue(ctx, transactionmanager.TxKey{}, tx)
	_, err = testRepo.GetByBookingInfo(txCtx, &object.BookingInfo{HotelID: booking.HotelID, RoomNumber: booking.RoomNumber, CheckIn: booking.CheckIn, CheckOut: booking.CheckOut, User: userdomain.User{ID: booking.UserID}})
	assert.Assert(t, errors.Is(err, error2.ErrBookingIsNotFound), fmt.Sprintf("error should be %v", error2.ErrBookingIsNotFound))
}

func TestBookingRepository_GetByHotel(t *testing.T) {
	booking := testBooking
	ctx := context.Background()
	tx, err := testDB.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("error starting transaction: %v", err)
	}

	defer tx.Rollback()

	txCtx := context.WithValue(ctx, transactionmanager.TxKey{}, tx)
	assert.Equal(t, nil, testRepo.Save(txCtx, &booking), fmt.Sprintf("testRepo.Save must be %v", nil))

	bookings, err := testRepo.GetByHotel(txCtx, booking.HotelID)
	if err != nil {
		t.Fatalf("error getting bookings: %v", err)
	}

	if bookings == nil || len(bookings) == 0 {
		t.Fatalf("bookings must not be nil and empty")
	}

	booking.ID = bookings[0].ID

	assert.Assert(t, equalBookings(booking, *bookings[0]), "booking must be equal to received booking")
}

func TestBookingRepository_GetByUser(t *testing.T) {
	booking := testBooking

	ctx := context.Background()
	tx, err := testDB.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("error starting transaction: %v", err)
	}

	defer tx.Rollback()

	txCtx := context.WithValue(ctx, transactionmanager.TxKey{}, tx)
	assert.Equal(t, nil, testRepo.Save(txCtx, &booking), fmt.Sprintf("testRepo.Save must be %v", nil))

	bookings, err := testRepo.GetByUser(txCtx, booking.UserID)
	if err != nil {
		t.Fatalf("error getting bookings: %v", err)
	}

	if bookings == nil || len(bookings) == 0 {
		t.Fatalf("bookings must not be nil and empty")
	}

	booking.ID = bookings[0].ID
	assert.Assert(t, equalBookings(booking, *bookings[0]), "booking must be equal to received booking")
}

func TestBookingRepository_GetBookingsByStatus(t *testing.T) {
	booking := testBooking
	booking.Status = bookingdomain.BookingStatusPaid
	ctx := context.Background()
	tx, err := testDB.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("error starting transaction: %v", err)
	}

	defer tx.Rollback()

	txCtx := context.WithValue(ctx, transactionmanager.TxKey{}, tx)
	assert.Equal(t, nil, testRepo.Save(txCtx, &booking), fmt.Sprintf("testRepo.Save must be %v", nil))
	bookings, err := testRepo.GetBookingsByStatus(txCtx, booking.Status)
	if err != nil {
		t.Fatalf("error getting bookings: %v", err)
	}

	if bookings == nil || len(bookings) == 0 {
		t.Fatalf("bookings must not be nil and empty")
	}

	booking.ID = bookings[0].ID
	assert.Assert(t, equalBookings(booking, *bookings[0]), "booking must be equal to received booking")
}
