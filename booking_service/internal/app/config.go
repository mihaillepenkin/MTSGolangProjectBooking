package app

import (
	"fmt"
	"os"

	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/app/http"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/app/jwt"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/infrastructure/config/grpcconfig"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/infrastructure/config/paymentconfig"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/infrastructure/config/postgresconfig"
	"gopkg.in/yaml.v3"
)

type Config struct {
	HTTPConfig     http.HTTPConfig               `yaml:"http"`
	PostgresConfig postgresconfig.PostgresConfig `yaml:"postgres"`
	GRPCConfig     grpcconfig.GRPCConfig         `yaml:"grpc"`
	PaymentConfig  paymentconfig.PaymentConfig   `yaml:"payment"`
	JWTConfig      jwt.JWTConfig                 `yaml:"jwt"`
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
