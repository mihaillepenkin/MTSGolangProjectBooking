package transactionmanager

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"gotest.tools/v3/assert"
)

var (
	testTxManager *TransactionManagerImpl[string]
	testDB        *sql.DB
)

func setup() error {
	container, err := postgres.Run(context.Background(), "postgres:15-alpine",
		postgres.WithDatabase("testDB-txManager"), postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		postgres.WithOrderedInitScripts("../../../../migrations/001_create_bookings_table.up.sql", "../../../../migrations/002_add_no_overlapping_constraint.up.sql"))

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
	testTxManager = NewTransactionManager[string](testDB)
	return nil
}

func TestMain(m *testing.M) {
	if err := setup(); err != nil {
		log.Fatal(err)
	}

	log.Println("Running tests...")
	code := m.Run()

	if err := testDB.Close(); err != nil {
		log.Fatal(err)
	}

	os.Exit(code)
}

func TestTransactionManagerImpl_ShouldReturnErrInTransaction(t *testing.T) {
	ctx := context.Background()
	_, err := testTxManager.InTransaction(ctx, func(ctx context.Context) (string, error) {
		tx, ok := GetTxFromCtx(ctx)
		assert.Assert(t, ok)
		_, err := tx.ExecContext(ctx, `INSERT INTO bookings (user_id, hotel_id, room_id, total_price, currency, check_in,
                             check_out, status, created_at, payment_id) VALUES 
                                                                ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`, uuid.New().String(), 1, 1, 100, "USD", time.Now().UTC(), time.Now().UTC(), "unpaid", time.Now(), "1")
		assert.Assert(t, err == nil)
		return "test", fmt.Errorf("test error")
	})

	assert.Assert(t, err != nil)

	rows, err := testDB.QueryContext(ctx, `SELECT * FROM bookings`)

	assert.Assert(t, err == nil)
	assert.Assert(t, !rows.Next())
}

func TestTransactionManagerImpl_WithCurTxInTransaction(t *testing.T) {
	ctx := context.Background()
	tx, err := testDB.BeginTx(ctx, nil)
	assert.Assert(t, err == nil)
	txCtx := context.WithValue(ctx, TxKey{}, tx)
	resp, err := testTxManager.InTransaction(txCtx, func(ctx context.Context) (string, error) {
		tx, ok := GetTxFromCtx(ctx)
		assert.Assert(t, ok)
		_, err := tx.ExecContext(ctx, `INSERT INTO bookings (user_id, hotel_id, room_id, total_price, currency, check_in,
                             check_out, status, created_at, payment_id) VALUES 
                                                                ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`, uuid.New().String(), 1, 1, 100, "USD", time.Now().UTC(), time.Now().UTC(), "unpaid", time.Now(), "1")
		assert.Assert(t, err == nil)
		return "test", nil
	})

	assert.Assert(t, err == nil && resp == "test")

	rows, err := tx.QueryContext(ctx, `SELECT * FROM bookings`)
	assert.Assert(t, err == nil)
	assert.Assert(t, rows.Next())
	err = tx.Rollback()
	assert.Assert(t, err == nil)
}

func TestTransactionManagerImpl_InTransaction(t *testing.T) {
	ctx := context.Background()
	resp, err := testTxManager.InTransaction(ctx, func(ctx context.Context) (string, error) {
		_, ok := GetTxFromCtx(ctx)
		assert.Assert(t, ok)
		return "test", nil
	})

	assert.Assert(t, err == nil && resp == "test")
}

func TestGetTxFromCtx(t *testing.T) {
	ctx := context.Background()
	_, ok := GetTxFromCtx(ctx)
	assert.Assert(t, !ok)

	tx, err := testDB.BeginTx(ctx, nil)
	assert.Assert(t, err == nil)
	ctx = context.WithValue(ctx, TxKey{}, tx)
	newTx, ok := GetTxFromCtx(ctx)
	assert.Assert(t, ok && tx == newTx)
}
