package booking

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/application/usecase/transactionmanager"
	bookingdomain "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking"
	error2 "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking/error"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking/object"
)

type BookingRepository struct {
	db *sql.DB
}

func NewBookingRepository(db *sql.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

func (b *BookingRepository) Save(ctx context.Context, booking *bookingdomain.Booking) error {
	var err error
	tx, ok := transactionmanager.GetTxFromCtx(ctx)
	if !ok {
		tx, err = b.db.BeginTx(ctx, nil)
		if err != nil {
			slog.Error("Error begin transaction", "error", err)
			return err
		}
		defer func() {
			if err != nil {
				if err2 := tx.Rollback(); err2 != nil {
					slog.Error("Error rollback", "error", err2)
				}
			}
		}()
	}

	query := `INSERT INTO bookings (user_id, hotel_name, room_number, total_price, currency, check_in,
                             check_out, status, created_at) VALUES 
                                                                ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
ON CONFLICT (user_id, hotel_name, room_number, check_in, check_out) DO UPDATE SET status = $8`

	result, execErr := tx.ExecContext(ctx, query, booking.UserID, booking.HotelName, booking.RoomNumber, booking.TotalPrice,
		booking.Currency, booking.CheckIn, booking.CheckOut, booking.Status, time.Now())
	if execErr != nil {
		slog.Error("Error saving booking", "error", execErr)
		err = execErr
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		slog.Error("Error getting rows affected", "error", err)
		return err
	}

	if rowsAffected == 0 {
		return error2.ErrBookingIsNotFound
	}

	if !ok {
		if commitErr := tx.Commit(); commitErr != nil {
			slog.Error("Error commit transaction", "error", commitErr)
			_ = tx.Rollback()
			return commitErr
		}
	}

	return nil
}

func (b *BookingRepository) Delete(ctx context.Context, booking *bookingdomain.Booking) error {
	var err error
	tx, ok := transactionmanager.GetTxFromCtx(ctx)
	if !ok {
		tx, err = b.db.BeginTx(ctx, nil)
		if err != nil {
			slog.Error("Error begin transaction", "error", err)
			return err
		}
		defer func() {
			if err != nil {
				if err2 := tx.Rollback(); err2 != nil {
					slog.Error("Error rollback", "error", err2)
				}
			}
		}()
	}

	query := `DELETE FROM bookings WHERE user_id = $1 AND hotel_name = $2 AND room_number = $3
              AND check_in = $4 AND check_out = $5`
	result, execErr := tx.ExecContext(ctx, query, booking.UserID, booking.HotelName, booking.RoomNumber, booking.CheckIn, booking.CheckOut)

	if execErr != nil {
		slog.Error("Error deleting booking", "error", execErr)
		err = execErr
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		slog.Error("Error getting rows affected", "error", err)
		return err
	}

	if rowsAffected == 0 {
		return error2.ErrBookingIsNotFound
	}

	if !ok {
		if commitErr := tx.Commit(); commitErr != nil {
			slog.Error("Error commit transaction", "error", commitErr)
			_ = tx.Rollback()
			return commitErr
		}
	}

	return nil
}

func (b *BookingRepository) GetDurationsByRoom(ctx context.Context, hotelName string, roomNumber string) ([][]time.Time, error) {
	var err error
	tx, err := b.db.BeginTx(ctx, nil)

	defer func() {
		if err != nil {
			if err2 := tx.Rollback(); err2 != nil {
				slog.Error("Error rollback", "error", err2)
			}
		}
	}()

	if err != nil {
		slog.Error("Error begin transaction", "error", err)
		return nil, err
	}

	query := `SELECT check_in, check_out
    FROM bookings
    WHERE hotel_name = $1 AND room_number = $2 and status ='paid'`

	rows, err := tx.QueryContext(ctx, query, hotelName, roomNumber)
	if err != nil {
		slog.Error("Error querying bookings", "error", err)
		return nil, err
	}

	defer rows.Close()

	durations := make([][]time.Time, 0)
	for rows.Next() {
		var checkIn, checkOut time.Time
		if err = rows.Scan(&checkIn, &checkOut); err != nil {
			slog.Error("Error scanning bookings", "error", err)
			return nil, err
		}

		durations = append(durations, []time.Time{checkIn, checkOut})
	}

	if err = rows.Err(); err != nil {
		slog.Error("Error scanning bookings", "error", err)
		return nil, err
	}

	if commitErr := tx.Commit(); commitErr != nil {
		_ = tx.Rollback()
		slog.Error("Error commit transaction", "error", commitErr)
		return nil, commitErr
	}

	return durations, nil
}

func (b *BookingRepository) IsIntersected(ctx context.Context, hotelName string, hotelRoom string, checkIn time.Time, checkOut time.Time) (bool, error) {
	var err error
	tx, ok := transactionmanager.GetTxFromCtx(ctx)
	if !ok {
		tx, err = b.db.BeginTx(ctx, nil)
		if err != nil {
			slog.Error("Error begin transaction", "error", err)
			return false, err
		}

		defer func() {
			if err != nil {
				if err2 := tx.Rollback(); err2 != nil {
					slog.Error("Error rollback", "error", err2)
				}
			}
		}()
	}

	query := `SELECT COUNT(*) FROM locked_rows
                WHERE hotel_name = $1 
                  AND hotel_number = $2
                  AND check_in < $4
                  AND check_out > $5`

	var count int
	execErr := tx.QueryRowContext(ctx, query, hotelName, hotelRoom, checkIn, checkOut).Scan(&count)
	if execErr != nil {
		slog.Error("Error checking booking", "error", execErr)
		return false, execErr
	}

	if !ok {
		if commitErr := tx.Commit(); commitErr != nil {
			_ = tx.Rollback()
			slog.Error("Error commit transaction", "error", commitErr)
			return false, commitErr
		}
	}

	return count > 0, nil
}

func (b *BookingRepository) GetByBookingInfo(ctx context.Context, bookingInfo *object.BookingInfo) (*bookingdomain.Booking, error) {
	var err error
	tx, ok := transactionmanager.GetTxFromCtx(ctx)
	if !ok {
		tx, err = b.db.BeginTx(ctx, nil)
		if err != nil {
			slog.Error("Error begin transaction", "error", err)
			return nil, err
		}

		defer func() {
			if err != nil {
				if err2 := tx.Rollback(); err2 != nil {
					slog.Error("Error rollback", "error", err2)
				}
			}
		}()
	}

	query := `SELECT id, user_id, hotel_name, room_number, total_price, currency, check_in,
                             check_out, status, created_at FROM bookings WHERE user_id = $1 AND hotel_name = $2 AND room_number = $3 AND check_in = $4 AND check_out = $5`

	var id string
	booking := &bookingdomain.Booking{}
	execErr := tx.QueryRowContext(ctx, query, bookingInfo.User.ID, bookingInfo.HotelName, bookingInfo.RoomNumber, bookingInfo.CheckIn, bookingInfo.CheckOut).Scan(&id, &booking.UserID, &booking.HotelName, &booking.RoomNumber, &booking.TotalPrice,
		&booking.Currency, &booking.CheckIn, &booking.CheckOut, &booking.Status)

	if errors.Is(execErr, sql.ErrNoRows) {
		return nil, error2.ErrBookingIsNotFound
	} else if execErr != nil {
		slog.Error("Error getting booking", "error", execErr)
		return nil, err
	}

	bookingID, idErr := object.NewBookingID(id)
	if idErr != nil {
		slog.Error("Error getting id", "error", idErr)
		return nil, err
	}

	booking.ID = bookingID

	if !ok {
		if commitErr := tx.Commit(); commitErr != nil {
			_ = tx.Rollback()
			slog.Error("Error commit transaction", "error", commitErr)
			return nil, commitErr
		}
	}

	return booking, nil
}

func (b *BookingRepository) GetByHotel(ctx context.Context, hotelName string) ([]*bookingdomain.Booking, error) {
	var err error
	tx, err := b.db.BeginTx(ctx, nil)

	defer func() {
		if err != nil {
			if err2 := tx.Rollback(); err2 != nil {
				slog.Error("Error rollback", "error", err2)
			}
		}
	}()

	if err != nil {
		slog.Error("Error begin transaction", "error", err)
		return nil, err
	}

	query := `SELECT id, user_id, hotel_name, room_number, total_price, currency, check_in, check_out, status FROM bookings 
WHERE hotel_name = $1`
	bookings := make([]*bookingdomain.Booking, 0)
	rows, err := tx.QueryContext(ctx, query, hotelName)

	if err != nil {
		slog.Error("Error getting booking", "error", err)
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var id string
		booking := &bookingdomain.Booking{}
		err = rows.Scan(&id, &booking.UserID, &booking.HotelName, &booking.RoomNumber, &booking.TotalPrice, &booking.Currency, &booking.CheckIn, &booking.CheckOut, &booking.Status)
		if err != nil {
			slog.Error("Error getting booking", "error", err)
			return nil, err
		}
		bookingID, idErr := object.NewBookingID(id)
		if idErr != nil {
			slog.Error("Error getting id", "error", idErr)
			err = idErr
			return nil, err
		}
		booking.ID = bookingID
		bookings = append(bookings, booking)
	}

	err = rows.Err()
	if err != nil {
		slog.Error("Error getting bookings", "error", err)
		return nil, err
	}

	if commitErr := tx.Commit(); commitErr != nil {
		_ = tx.Rollback()
		slog.Error("Error commit transaction", "error", commitErr)
		return nil, commitErr
	}

	return bookings, nil
}

func (b *BookingRepository) GetByUser(ctx context.Context, userID string) ([]*bookingdomain.Booking, error) {
	var err error
	tx, err := b.db.BeginTx(ctx, nil)
	if err != nil {
		slog.Error("Error begin transaction", "error", err)
		return nil, err
	}

	defer func() {
		if err != nil {
			if err2 := tx.Rollback(); err2 != nil {
				slog.Error("Error rollback", "error", err2)
			}
		}
	}()

	query := `SELECT id, user_id, hotel_name, room_number, total_price, currency, check_in, check_out, status FROM bookings
WHERE user_id = $1`
	bookings := make([]*bookingdomain.Booking, 0)
	rows, err := tx.QueryContext(ctx, query, userID)
	if err != nil {
		slog.Error("Error getting booking", "error", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var id string
		booking := &bookingdomain.Booking{}
		err = rows.Scan(ctx, query, &id, &booking.UserID, &booking.HotelName, &booking.RoomNumber, &booking.TotalPrice, &booking.Currency, &booking.CheckIn, &booking.CheckOut, &booking.Status)
		if err != nil {
			slog.Error("Error getting booking", "error", err)
			return nil, err
		}

		bookingID, idErr := object.NewBookingID(id)
		if idErr != nil {
			slog.Error("Error getting id", "error", idErr)
			err = idErr
			return nil, err
		}

		booking.ID = bookingID

		bookings = append(bookings, booking)
	}

	err = rows.Err()
	if err != nil {
		slog.Error("Error getting bookings", "error", err)
		return nil, err
	}

	if commitErr := tx.Commit(); commitErr != nil {
		_ = tx.Rollback()
		slog.Error("Error commit transaction", "error", commitErr)
		return nil, commitErr
	}

	return bookings, nil
}
