package app

import (
	"database/sql"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var (
	ConfigPath = "config.yml"
)

type App struct {
	db     *sql.DB
	config *Config
}

func NewApp() (*App, error) {
	setLogger()
	cfg, err := LoadConfig(ConfigPath)
	if err != nil {
		slog.Error("Error loading config: ", "error", err)
		return nil, err
	}

	db, err := sql.Open("postgres", cfg.PostgresConfig.DSN)
	if err != nil {
		slog.Error("Error opening connection to database: ", "error", err)
		return nil, err
	}

	db.SetMaxIdleConns(cfg.PostgresConfig.MaxIdleConnections)
	db.SetMaxOpenConns(cfg.PostgresConfig.MaxConnections)
	db.SetConnMaxLifetime(cfg.PostgresConfig.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.PostgresConfig.ConnMaxIdleTime)

	if err := db.Ping(); err != nil {
		db.Close()
		slog.Error("Error connecting to database: ", "error", err)
		return nil, err
	}

	slog.Info("Successfully connected to PostgreSQL")
	return &App{db: db, config: cfg}, nil
}

func (app *App) Start() error {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	slog.Info("Starting application...")

	err := app.RunMigrations()
	if err != nil {
		slog.Error("Error starting application: ", "error", err)
		app.db.Close()
		return err
	}

	<-signalChan
	slog.Info("Shutting down...")

	dbErr := app.db.Close()
	if dbErr != nil {
		slog.Error("Error closing database: ", "error", dbErr)
		return dbErr
	}

	slog.Info("Application stopped successfully")
	return nil
}

func (a *App) RunMigrations() error {
	driver, err := postgres.WithInstance(a.db, &postgres.Config{})
	if err != nil {
		slog.Error("Error creating database driver", "error", err)
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		slog.Error("Error creating migration instance", "error", err)
		return err
	}
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		slog.Error("Error running migration", "error", err)
		return err
	}

	version, dirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		slog.Error("Error getting version", "error", err)
		return err
	}

	slog.Info("Migrations complete", "version", version, "dirty", dirty)
	return nil
}

func (a *App) ForceMigration(versionToMigrate int) error {
	driver, err := postgres.WithInstance(a.db, &postgres.Config{})
	if err != nil {
		slog.Error("Error creating database driver", "error", err)
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		slog.Error("Error creating migration instance", "error", err)
		return err
	}

	err = m.Force(versionToMigrate)
	if err != nil {
		slog.Error("Error running migration", "error", err)
		return err
	}

	version, dirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		slog.Error("Error getting version", "error", err)
		return err
	}
	slog.Info("Migration complete", "version", version, "dirty", dirty)
	return nil
}

func setLogger() {
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
