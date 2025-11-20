package main

import (
	"auth_service/internal/application/usecase"
	"auth_service/internal/config"
	"auth_service/internal/infrastructure/handlers"
	"auth_service/internal/infrastructure/postgres"
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	// Initialize repositories
	userRepo := &postgres.UserRepoPG{DB: db}
	jwtRepo := &postgres.JWTRepoPG{DB: db}

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