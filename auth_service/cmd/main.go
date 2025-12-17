package main

import (
	"auth_service/internal/application/usecase"
	"auth_service/internal/config"
	"auth_service/internal/infrastructure/handlers"
	authpostgres "auth_service/internal/infrastructure/postgres"
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to PostgreSQL
	if cfg.Port == "" {
		cfg.Port = "8081"
	}
	dsn := "host=" + cfg.DBHost + " port=" + cfg.DBPort +
		" user=" + cfg.DBUser + " password=" + cfg.DBPassword +
		" dbname=" + cfg.DBName + " sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to Postgres: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping Postgres: %v", err)
	}

	if err := runMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	userRepo := &authpostgres.UserRepoPG{DB: db}
	jwtRepo := &authpostgres.JWTRepoPG{DB: db}

	// Initialize JWT service
	jwtService := usecase.JWTService(cfg.JWTSecret)

	// Initialize interfaces layer (formerly usecase)
	userService := usecase.UserService(userRepo, jwtService, jwtRepo)

	// Initialize HTTP handlers
	authHandler := handlers.NewHandler(userService)

	// Setup HTTP router
	routers := authHandler.SetupRoutes()

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           routers,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// Run server in goroutine
	go func() {
		log.Printf("Auth Service running on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Println("Shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}

func runMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Printf("Error creating database driver, error: %v", err)
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		log.Printf("Error creating migration instance, error: %v", err)
		return err
	}
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Printf("Error running migration, error: %v", err)
		return err
	}

	version, dirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		log.Printf("Error getting version, error: %v", err)
		return err
	}

	log.Printf("Migrations complete, version %v, dirty %v", version, dirty)
	return nil
}
