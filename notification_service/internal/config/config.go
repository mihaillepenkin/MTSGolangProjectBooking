package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)
type KafkaConfig struct {
	Brokers   []string
	Topic     string
	GroupID   string
}

type SMTPConfig struct {
	FromEmail string
	Password  string
	SMTPHost  string
	SMTPPort  string
}

type Config struct {
	Kafka KafkaConfig
	SMTP  SMTPConfig
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: .env file not found, using environment variables")
	}

	getEnv := func(key, desc string) (string, error) {
		value := os.Getenv(key)
		if value == "" {
			return "", fmt.Errorf("missing required environment variable: %s (%s)", key, desc)
		}
		return value, nil
	}

	// Kafka
	kafkaBrokersRaw, err := getEnv("KAFKA_BROKERS", "comma-separated list of Kafka brokers, e.g., 'localhost:9092'")
	if err != nil {
		return nil, err
	}
	kafkaBrokers := splitAndTrim(kafkaBrokersRaw, ",")

	kafkaTopic, err := getEnv("KAFKA_TOPIC", "Kafka topic name")
	if err != nil {
		return nil, err
	}

	kafkaGroupID, err := getEnv("KAFKA_GROUP_ID", "Kafka consumer group ID")
	if err != nil {
		return nil, err
	}

	// SMTP
	fromEmail, err := getEnv("NOTIFICATION_FROM_EMAIL", "sender email address")
	if err != nil {
		return nil, err
	}

	smtpPassword, err := getEnv("NOTIFICATION_EMAIL_PASSWORD", "SMTP/app password for sender email")
	if err != nil {
		return nil, err
	}

	smtpHost, err := getEnv("NOTIFICATION_SMTP_HOST", "SMTP server host, e.g., smtp.gmail.com")
	if err != nil {
		return nil, err
	}

	smtpPort, err := getEnv("NOTIFICATION_SMTP_PORT", "SMTP server port, e.g., 587")
	if err != nil {
		return nil, err
	}

	return &Config{
		Kafka: KafkaConfig{
			Brokers: kafkaBrokers,
			Topic:   kafkaTopic,
			GroupID: kafkaGroupID,
		},
		SMTP: SMTPConfig{
			FromEmail: fromEmail,
			Password:  smtpPassword,
			SMTPHost:  smtpHost,
			SMTPPort:  smtpPort,
		},
	}, nil
}

func splitAndTrim(s, sep string) []string {
	parts := strings.Split(s, sep)
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}