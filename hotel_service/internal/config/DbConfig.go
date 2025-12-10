package config

import (
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type PostgresConfig struct {
	Dsn                string        `yaml:"dsn"`
	MaxConnections     int           `yaml:"max_connections"`
	MaxIdleConnections int           `yaml:"max_idle_connections"`
	ConnMaxLifetime    time.Duration `yaml:"conn_max_lifetime"`
	ConnMaxIdleTime    time.Duration `yaml:"conn_max_idle_time"`
}
