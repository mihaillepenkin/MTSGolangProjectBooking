package config

import (
	"database/sql"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type DbConfig struct {
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassword string
	DbName     string
}

func ConfigureDb() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		slog.Error(err.Error())
		return nil
	}

	cfg := DbConfig{
		DbHost:     os.Getenv("DB_HOST"),
		DbPort:     os.Getenv("DB_PORT"),
		DbUser:     os.Getenv("DB_USER"),
		DbPassword: os.Getenv("DB_PASSWORD"),
		DbName:     os.Getenv("DB_NAME"),
	}

	dsn := "host=" + cfg.DbHost + " port=" + cfg.DbPort + " user=" + cfg.DbUser + " password=" + cfg.DbPassword + " dbname=" + cfg.DbName + " sslmode=disable"

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

	return db
}
