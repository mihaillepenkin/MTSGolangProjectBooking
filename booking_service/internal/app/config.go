package app

import (
	"fmt"
	"os"

	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/infrastructure/config/grpcconfig"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/infrastructure/config/postgresconfig"
	"gopkg.in/yaml.v3"
)

type Config struct {
	HTTPConfig     HTTPConfig                    `yaml:"http"`
	PostgresConfig postgresconfig.PostgresConfig `yaml:"postgres"`
	GRPCConfig     grpcconfig.GRPCConfig         `yaml:"grpc"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %v", err)
	}
	return &config, nil
}
