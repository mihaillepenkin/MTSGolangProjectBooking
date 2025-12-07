package httpconfig

import "time"

type HTTPConfig struct {
	Address        string        `yaml:"address"`
	Host           string        `yaml:"host"`
	Port           string        `yaml:"port"`
	AllowedOrigins []string      `yaml:"allowed_origins"`
	ReadTimeout    time.Duration `yaml:"read_timeout"`
}
