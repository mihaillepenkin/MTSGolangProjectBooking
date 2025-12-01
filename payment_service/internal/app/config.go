package app

import (
	"fmt"
	"os"

	"github.com/mihaillepenkin/MTSGolangProjectBooking/payment_service/internal/adapter/httpconfig"
	"gopkg.in/yaml.v3"
)

type Config struct {
	HTTPConfig httpconfig.HTTPConfig `yaml:"http"`
	SecretKey  string                `yaml:"secret_key"`
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
