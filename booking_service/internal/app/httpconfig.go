package app

import "time"

type HTTPConfig struct {
	Address      string        `yaml:"address"`
	Host         string        `yaml:"host"`
	Port         string        `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	AllowOrigins []string      `yaml:"allowed_origins"`
}
