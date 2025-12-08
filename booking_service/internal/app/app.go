package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/usecase/case/transactionmanager"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	ConfigPath = "config.yml"
)

type App struct {
	db           *sql.DB
	config       *Config
	server       *http.Server
	handlers     *Handlers
	services     *Services
	repositories *Repositories
	grpcConn     *grpc.ClientConn
	kafkaWriter  *kafka.Writer
	logger       *slog.Logger
}

func NewApp() (*App, error) {
	setLogger()
	cfg, err := LoadConfig(ConfigPath)
	logger := slog.Default().With("component", "app")
	if err != nil {
		logger.Error("Error loading config: ", "error", err)
		return nil, err
	}

	db, err := sql.Open("postgres", cfg.PostgresConfig.DSN)
	if err != nil {
		logger.Error("Error opening connection to database: ", "error", err)
		return nil, err
	}

	db.SetMaxIdleConns(cfg.PostgresConfig.MaxIdleConnections)
	db.SetMaxOpenConns(cfg.PostgresConfig.MaxConnections)
	db.SetConnMaxLifetime(cfg.PostgresConfig.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.PostgresConfig.ConnMaxIdleTime)

	if err := db.Ping(); err != nil {
		db.Close()
		logger.Error("Error connecting to database: ", "error", err)
		return nil, err
	}

	conn, err := grpc.NewClient(cfg.GRPCConfig.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		db.Close()
		logger.Error("Error creating grpc client: ", "error", err)
		return nil, err
	}

	kafkaWriter := &kafka.Writer{Addr: kafka.TCP(cfg.KafkaConfig.Address), Topic: cfg.KafkaConfig.Topic, Balancer: &kafka.RoundRobin{}}
	txManager := transactionmanager.NewTransactionManager[string](db)
	repos := NewRepositories(db, cfg.PaymentConfig, conn, kafkaWriter)
	services := NewServices(cfg, repos, txManager)
	handlers := NewHandlers(services)
	handler := handlers.registerRoutes(cfg)

	logger.Info("Successfully connected to PostgreSQL")
	return &App{db: db, config: cfg, server: &http.Server{Addr: cfg.HTTPConfig.Address, Handler: handler, WriteTimeout: cfg.HTTPConfig.WriteTimeout, ReadTimeout: cfg.HTTPConfig.ReadTimeout},
		repositories: repos, services: services, handlers: handlers, grpcConn: conn, kafkaWriter: kafkaWriter, logger: logger}, nil
}

func (app *App) Run() error {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	app.logger.Info("Starting application...")

	err := app.RunMigrations()
	if err != nil {
		app.logger.Error("Error starting application: ", "error", err)
		app.closeAll()
		return err
	}

	go func() {
		app.logger.Info(fmt.Sprintf("Starting server on port %s", app.config.HTTPConfig.Port))
		if err := app.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			app.logger.Error("Server failed", "error", err)
			app.closeAll()
			os.Exit(1)
		}
	}()

	<-signalChan
	app.logger.Info("Shutting down...")

	app.logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	servErr := app.server.Shutdown(ctx)
	if servErr != nil {
		app.logger.Error("Error shutting down server", "error", servErr)
	}

	app.logger.Info("Shutting down db...")
	dbErr := app.db.Close()
	if dbErr != nil {
		app.logger.Error("Error closing database: ", "error", dbErr)
	}

	app.logger.Info("Shutting down grpc connection...")

	grpcErr := app.grpcConn.Close()
	if grpcErr != nil {
		app.logger.Error("Error closing grpc connection: ", "error", grpcErr)
	}

	kafkaErr := app.kafkaWriter.Close()
	if kafkaErr != nil {
		app.logger.Error("Error closing kafka: ", "error", kafkaErr)
	}

	if servErr != nil || dbErr != nil || grpcErr != nil || kafkaErr != nil {
		return fmt.Errorf("error during application shutdown: %v, %v, %v, %v", servErr, dbErr, grpcErr, kafkaErr)
	}

	app.logger.Info("Application stopped successfully")
	return nil
}

func (app *App) RunMigrations() error {
	driver, err := postgres.WithInstance(app.db, &postgres.Config{})
	if err != nil {
		app.logger.Error("Error creating database driver", "error", err)
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		app.logger.Error("Error creating migration instance", "error", err)
		return err
	}
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		app.logger.Error("Error running migration", "error", err)
		return err
	}

	version, dirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		app.logger.Error("Error getting version", "error", err)
		return err
	}

	app.logger.Info("Migrations complete", "version", version, "dirty", dirty)
	return nil
}

func (app *App) ForceMigration(versionToMigrate int) error {
	driver, err := postgres.WithInstance(app.db, &postgres.Config{})
	if err != nil {
		app.logger.Error("Error creating database driver", "error", err)
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		app.logger.Error("Error creating migration instance", "error", err)
		return err
	}

	err = m.Force(versionToMigrate)
	if err != nil {
		app.logger.Error("Error running migration", "error", err)
		return err
	}

	version, dirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		app.logger.Error("Error getting version", "error", err)
		return err
	}
	app.logger.Info("Migration complete", "version", version, "dirty", dirty)
	return nil
}

func (app *App) closeAll() {
	err := app.db.Close()
	if err != nil {
		app.logger.Error("Error closing database connection: ", "error", err)
	}
	err = app.grpcConn.Close()
	if err != nil {
		app.logger.Error("Error closing gRPC connection: ", "error", err)
	}
	err = app.kafkaWriter.Close()
	if err != nil {
		app.logger.Error("Error closing kafka writer: ", "error", err)
	}
}

func setLogger() {
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
