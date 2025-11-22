package postgresconfig

import (
	"time"
)

type PostgresConfig struct {
	DSN                string        `yaml:"dsn"`
	MaxConnections     int           `yaml:"max_connections"`
	MaxIdleConnections int           `yaml:"max_idle_connections"`
	ConnMaxLifetime    time.Duration `yaml:"conn_max_lifetime"`
	ConnMaxIdleTime    time.Duration `yaml:"conn_max_idle_time"`
}
