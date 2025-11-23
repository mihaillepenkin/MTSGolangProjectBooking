package config

import (
	"database/sql"
	"os"
)

type DbConfig struct {
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassword string
	DbName     string
}

func ConfigureDb() (*sql.DB, error) {
	cfg := DbConfig{
		DbHost:     os.Getenv("DB_HOST"),
		DbPort:     os.Getenv("DB_PORT"),
		DbUser:     os.Getenv("DB_USER"),
		DbPassword: os.Getenv("DB_PASSWORD"),
		DbName:     os.Getenv("DB_NAME"),
	}

	dsn := "host=" + cfg.DbHost + " port=" + cfg.DbPort + " user=" + cfg.DbUser + " password=" + cfg.DbPassword + " dbname=" + cfg.DbName + " sslmode=disable"

	return sql.Open("postgres", dsn)
}
