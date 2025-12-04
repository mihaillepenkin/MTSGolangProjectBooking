package httpconfig

import "time"

type HTTPConfig struct {
	Address                string        `yaml:"address"`
	Host                   string        `yaml:"host"`
	Port                   string        `yaml:"port"`
	ReadTimeout            time.Duration `yaml:"read_timeout"`
	WriteTimeout           time.Duration `yaml:"write_timeout"`
	AllowedOrigins         []string      `yaml:"allowed_origins"`
	WebhookHandlerEndpoint string        `yaml:"webhook_handler_endpoint"`
}
