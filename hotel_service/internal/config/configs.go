package config

import (
	"database/sql"
	"errors"
	grpc2 "hotel_service/internal/infrastructure/grpc"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/Vlad-Ali/MTSGolangProjectBooking-protos/gen/proto/hotel"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"
)

type Config struct {
	PostgresConfig PostgresConfig `yaml:"postgres"`
	GrpcConfig     GrpcConfig     `yaml:"grpc"`
	HttpConfig     HttpConfig     `yaml:"http"`
}

func LoadConfig() *Config {
	data, err := os.ReadFile("config.yml")
	if err != nil {
		slog.Error(err.Error())
		return nil
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		slog.Error(err.Error())
		return nil
	}

	slog.Info("Successfully loaded config")
	return &config
}

func (cfg *Config) ConnectToDb() *sql.DB {
	db, err := sql.Open("postgres", cfg.PostgresConfig.Dsn)
	if err != nil {
		slog.Error(err.Error(), "DSN", cfg.PostgresConfig.Dsn)
		return nil
	}
	/*defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			slog.Error(err.Error(), "DSN", cfg.PostgresConfig.Dsn)
		}
	}(db)*/

	err = db.Ping()
	if err != nil {
		slog.Error(err.Error(), "DSN", cfg.PostgresConfig.Dsn)
		return nil
	} else {
		slog.Info("Successfully connected to database")
	}

	return db
}

func (cfg *Config) RunMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		slog.Error(err.Error())
		return err
	}

	version, dirty, err := m.Version()
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	slog.Info("Migrations complete", "version", version, "dirty", dirty)
	return nil
}

func (cfg *Config) ConfigureGrpcServer(hotelGrpcService *grpc2.HotelGrpcService) (*grpc.Server, net.Listener) {
	grpcServer := grpc.NewServer()
	hotel.RegisterHotelServer(grpcServer, hotelGrpcService)

	listener, err := net.Listen("tcp", cfg.GrpcConfig.Address)
	if err != nil {
		slog.Error(err.Error())
		return nil, nil
	}

	slog.Info("Successfully configured gRPC server")
	return grpcServer, listener
}

func (cfg *Config) ConfigureHttpServer(handler *http.Handler) *http.Server {
	server := &http.Server{
		Addr:         cfg.HttpConfig.Address,
		Handler:      *handler,
		ReadTimeout:  cfg.HttpConfig.ReadTimeout,
		WriteTimeout: cfg.HttpConfig.WriteTimeout,
	}

	slog.Info("Successfully configured http server")
	return server
}
