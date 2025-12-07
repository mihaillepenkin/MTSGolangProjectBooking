package config

import (
	"database/sql"
	"log/slog"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func ConfigureDb() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		slog.Error(err.Error())
		return nil
	}

	cfg := &PostgresConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}

	dsn := "host=" + cfg.Host + " port=" + cfg.Port + " user=" + cfg.User + " password=" + cfg.Password + " dbname=" + cfg.Name + " sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		slog.Error(err.Error())
		return nil
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			slog.Error(err.Error())
		}
	}(db)

	err = db.Ping()
	if err != nil {
		slog.Error(err.Error())
		return nil
	} else {
		slog.Info("Successfully connected to database")
	}

	return db
}

func RunMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	err = m.Up()
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	version, dirty, err := m.Version()
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	slog.Info("Migrations complete", "version", version, "dirty", dirty)
	return nil
}
