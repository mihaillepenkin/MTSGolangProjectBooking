package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	bookingdomain "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking"
	error2 "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking/error"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking/object"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/usecase/case/transactionmanager"
)

type BookingRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewBookingRepository(db *sql.DB) *BookingRepository {
	return &BookingRepository{db: db, logger: slog.Default().With("component", "booking_repository")}
}

func (b *BookingRepository) Save(ctx context.Context, booking *bookingdomain.Booking) error {
	var err error
	tx, ok := transactionmanager.GetTxFromCtx(ctx)
	if !ok {
		tx, err = b.db.BeginTx(ctx, nil)
		if err != nil {
			b.logger.Error("Error begin transaction", "error", err)
			return err
		}
		defer func() {
			if err != nil {
				if err2 := tx.Rollback(); err2 != nil {
					b.logger.Error("Error rollback", "error", err2)
				}
			}
		}()
	}

	query := `INSERT INTO bookings (user_id, hotel_id, room_id, total_price, currency, check_in,
                             check_out, status, created_at, payment_id) VALUES 
                                                                ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) 
ON CONFLICT (user_id, hotel_id, room_id, check_in, check_out) DO UPDATE SET status = $8`

	result, execErr := tx.ExecContext(ctx, query, booking.UserID, booking.HotelID, booking.RoomID, booking.TotalPrice,
		booking.Currency, booking.CheckIn, booking.CheckOut, booking.Status, time.Now(), booking.PaymentID)
	if execErr != nil {
		b.logger.Error("Error saving postgres", "error", execErr)
		err = execErr
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		b.logger.Error("Error getting rows affected", "error", err)
		return err
	}

	if rowsAffected == 0 {
		return error2.ErrBookingIsNotFound
	}

	if !ok {
		if commitErr := tx.Commit(); commitErr != nil {
			b.logger.Error("Error commit transaction", "error", commitErr)
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
			b.logger.Error("Error begin transaction", "error", err)
			return err
		}
		defer func() {
			if err != nil {
				if err2 := tx.Rollback(); err2 != nil {
					b.logger.Error("Error rollback", "error", err2)
				}
			}
		}()
	}

	query := `DELETE FROM bookings WHERE user_id = $1 AND hotel_id = $2 AND room_id = $3
              AND check_in = $4 AND check_out = $5`
	result, execErr := tx.ExecContext(ctx, query, booking.UserID, booking.HotelID, booking.RoomID, booking.CheckIn, booking.CheckOut)

	if execErr != nil {
		b.logger.Error("Error deleting postgres", "error", execErr)
		err = execErr
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		b.logger.Error("Error getting rows affected", "error", err)
		return err
	}

	if rowsAffected == 0 {
		return error2.ErrBookingIsNotFound
	}

	if !ok {
		if commitErr := tx.Commit(); commitErr != nil {
			b.logger.Error("Error commit transaction", "error", commitErr)
			_ = tx.Rollback()
			return commitErr
		}
	}

	return nil
}

func (b *BookingRepository) GetDurationsByRoom(ctx context.Context, hotelID int64, roomID int64) ([][]time.Time, error) {
	var err error
	tx, ok := transactionmanager.GetTxFromCtx(ctx)
	if !ok {
		tx, err = b.db.BeginTx(ctx, nil)
		if err != nil {
			b.logger.Error("Error begin transaction", "error", err)
			return nil, err
		}

		defer func() {
			if err != nil {
				if err2 := tx.Rollback(); err2 != nil {
					b.logger.Error("Error rollback", "error", err2)
				}
			}
		}()
	}

	query := `SELECT check_in, check_out
    FROM bookings
    WHERE hotel_id = $1 AND room_id = $2 and status ='paid'`

	rows, err := tx.QueryContext(ctx, query, hotelID, roomID)
	if err != nil {
		b.logger.Error("Error querying bookings", "error", err)
		return nil, err
	}

	defer rows.Close()

	durations := make([][]time.Time, 0)
	for rows.Next() {
		var checkIn, checkOut time.Time
		if err = rows.Scan(&checkIn, &checkOut); err != nil {
			b.logger.Error("Error scanning bookings", "error", err)
			return nil, err
		}

		durations = append(durations, []time.Time{checkIn, checkOut})
	}

	if err = rows.Err(); err != nil {
		b.logger.Error("Error scanning bookings", "error", err)
		return nil, err
	}

	if !ok {
		if commitErr := tx.Commit(); commitErr != nil {
			b.logger.Error("Error commit transaction", "error", commitErr)
			_ = tx.Rollback()
			return nil, commitErr
		}
	}
	return durations, nil
}

func (b *BookingRepository) IsIntersected(ctx context.Context, hotelID int64, roomID int64, checkIn time.Time, checkOut time.Time) (bool, error) {
	var err error
	tx, ok := transactionmanager.GetTxFromCtx(ctx)
	if !ok {
		tx, err = b.db.BeginTx(ctx, nil)
		if err != nil {
			b.logger.Error("Error begin transaction", "error", err)
			return false, err
		}

		defer func() {
			if err != nil {
				if err2 := tx.Rollback(); err2 != nil {
					b.logger.Error("Error rollback", "error", err2)
				}
			}
		}()
	}

	query := `SELECT COUNT(*) FROM bookings
                WHERE hotel_id = $1 
                  AND room_id = $2
                  AND check_in < $4
                  AND check_out > $3`

	var count int
	execErr := tx.QueryRowContext(ctx, query, hotelID, roomID, checkIn, checkOut).Scan(&count)
	if execErr != nil {
		b.logger.Error("Error checking postgres", "error", execErr)
		return false, execErr
	}

	if !ok {
		if commitErr := tx.Commit(); commitErr != nil {
			_ = tx.Rollback()
			b.logger.Error("Error commit transaction", "error", commitErr)
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
			b.logger.Error("Error begin transaction", "error", err)
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

	query := `SELECT id, user_id, hotel_id, room_id, total_price, currency, check_in,
                             check_out, status, payment_id FROM bookings WHERE user_id = $1 AND hotel_id = $2 AND room_id = $3 AND check_in = $4 AND check_out = $5`

	var id string
	booking := &bookingdomain.Booking{}
	execErr := tx.QueryRowContext(ctx, query, bookingInfo.User.ID, bookingInfo.HotelID, bookingInfo.RoomID, bookingInfo.CheckIn, bookingInfo.CheckOut).Scan(&id, &booking.UserID, &booking.HotelID, &booking.RoomID, &booking.TotalPrice,
		&booking.Currency, &booking.CheckIn, &booking.CheckOut, &booking.Status, &booking.PaymentID)

	if errors.Is(execErr, sql.ErrNoRows) {
		return nil, error2.ErrBookingIsNotFound
	} else if execErr != nil {
		b.logger.Error("Error getting postgres", "error", execErr)
		return nil, err
	}

	bookingID, idErr := object.NewBookingID(id)
	if idErr != nil {
		b.logger.Error("Error getting id", "error", idErr)
		return nil, err
	}

	booking.ID = bookingID

	if !ok {
		if commitErr := tx.Commit(); commitErr != nil {
			_ = tx.Rollback()
			b.logger.Error("Error commit transaction", "error", commitErr)
			return nil, commitErr
		}
	}

	return booking, nil
}

func (b *BookingRepository) GetByHotel(ctx context.Context, hotelID int64) ([]*bookingdomain.Booking, error) {
	var err error
	tx, ok := transactionmanager.GetTxFromCtx(ctx)
	if !ok {
		tx, err = b.db.BeginTx(ctx, nil)
		if err != nil {
			b.logger.Error("Error begin transaction", "error", err)
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

	query := `SELECT id, user_id, hotel_id, room_id, total_price, currency, check_in, check_out, status, payment_id FROM bookings 
WHERE hotel_id = $1`
	bookings := make([]*bookingdomain.Booking, 0)
	rows, err := tx.QueryContext(ctx, query, hotelID)

	if err != nil {
		b.logger.Error("Error getting postgres", "error", err)
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var id string
		booking := &bookingdomain.Booking{}
		err = rows.Scan(&id, &booking.UserID, &booking.HotelID, &booking.RoomID, &booking.TotalPrice, &booking.Currency, &booking.CheckIn, &booking.CheckOut, &booking.Status, &booking.PaymentID)
		if err != nil {
			b.logger.Error("Error getting postgres", "error", err)
			return nil, err
		}
		bookingID, idErr := object.NewBookingID(id)
		if idErr != nil {
			b.logger.Error("Error getting id", "error", idErr)
			err = idErr
			return nil, err
		}
		booking.ID = bookingID
		bookings = append(bookings, booking)
	}

	err = rows.Err()
	if err != nil {
		b.logger.Error("Error getting bookings", "error", err)
		return nil, err
	}

	if !ok {
		if commitErr := tx.Commit(); commitErr != nil {
			_ = tx.Rollback()
			b.logger.Error("Error commit transaction", "error", commitErr)
			return nil, commitErr
		}
	}

	return bookings, nil
}

func (b *BookingRepository) GetByUser(ctx context.Context, userID string) ([]*bookingdomain.Booking, error) {
	var err error
	tx, ok := transactionmanager.GetTxFromCtx(ctx)
	if !ok {
		tx, err = b.db.BeginTx(ctx, nil)
		if err != nil {
			b.logger.Error("Error begin transaction", "error", err)
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

	query := `SELECT id, user_id, hotel_id, room_id, total_price, currency, check_in, check_out, status, payment_id FROM bookings
WHERE user_id = $1`
	bookings := make([]*bookingdomain.Booking, 0)
	rows, err := tx.QueryContext(ctx, query, userID)
	if err != nil {
		b.logger.Error("Error getting postgres", "error", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var id string
		booking := &bookingdomain.Booking{}
		err = rows.Scan(&id, &booking.UserID, &booking.HotelID, &booking.RoomID, &booking.TotalPrice, &booking.Currency, &booking.CheckIn, &booking.CheckOut, &booking.Status, &booking.PaymentID)
		if err != nil {
			b.logger.Error("Error getting postgres", "error", err)
			return nil, err
		}

		bookingID, idErr := object.NewBookingID(id)
		if idErr != nil {
			b.logger.Error("Error getting id", "error", idErr)
			err = idErr
			return nil, err
		}

		booking.ID = bookingID

		bookings = append(bookings, booking)
	}

	err = rows.Err()
	if err != nil {
		b.logger.Error("Error getting bookings", "error", err)
		return nil, err
	}

	if !ok {
		if commitErr := tx.Commit(); commitErr != nil {
			_ = tx.Rollback()
			b.logger.Error("Error commit transaction", "error", commitErr)
			return nil, commitErr
		}
	}

	return bookings, nil
}

func (b *BookingRepository) GetBookingsByStatus(ctx context.Context, status bookingdomain.BookingStatus) ([]*bookingdomain.Booking, error) {
	var err error
	tx, ok := transactionmanager.GetTxFromCtx(ctx)
	if !ok {
		tx, err = b.db.BeginTx(ctx, nil)
		if err != nil {
			b.logger.Error("Error begin transaction", "error", err)
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

	query := `SELECT id, user_id, hotel_id, room_id, total_price, currency, check_in, check_out, status, payment_id FROM bookings
WHERE status = $1 AND NOW() - created_at < INTERVAL '6 hour'`

	bookings := make([]*bookingdomain.Booking, 0)
	rows, err := tx.QueryContext(ctx, query, string(status))
	if err != nil {
		b.logger.Error("Error getting bookings", "error", err)
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var id string
		booking := &bookingdomain.Booking{}
		err = rows.Scan(&id, &booking.UserID, &booking.HotelID, &booking.RoomID, &booking.TotalPrice, &booking.Currency, &booking.CheckIn, &booking.CheckOut, &booking.Status, &booking.PaymentID)
		if err != nil {
			b.logger.Error("Error getting postgres", "error", err)
			return nil, err
		}

		bookingID, idErr := object.NewBookingID(id)
		if idErr != nil {
			b.logger.Error("Error getting id", "error", idErr)
			err = idErr
			return nil, err
		}

		booking.ID = bookingID

		bookings = append(bookings, booking)
	}

	err = rows.Err()
	if err != nil {
		b.logger.Error("Error getting bookings", "error", err)
		return nil, err
	}

	if !ok {
		if commitErr := tx.Commit(); commitErr != nil {
			_ = tx.Rollback()
			b.logger.Error("Error commit transaction", "error", commitErr)
			return nil, commitErr
		}
	}

	return bookings, nil
}
