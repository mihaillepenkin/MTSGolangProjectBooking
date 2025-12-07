package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	ConfigPath = "config.yml"
)

type App struct {
	handlers *Handlers
	services *Services
	server   *http.Server
	config   *Config
}

func NewApp() (*App, error) {
	setLogger()
	cfg, err := LoadConfig(ConfigPath)
	if err != nil {
		slog.Error("Error loading config: ", "error", err)
		return nil, err
	}

	services := NewServices(cfg.SecretKey)
	handlers := NewHandlers(cfg.HTTPConfig, services)
	handler := handlers.registerRoutes(cfg)
	slog.Info("Successfully initialized app")
	return &App{services: services, handlers: handlers, config: cfg,
		server: &http.Server{Addr: cfg.HTTPConfig.Address,
			Handler: handler, ReadTimeout: cfg.HTTPConfig.ReadTimeout}}, nil
}

func (app *App) Start() error {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	slog.Info("Starting application...")

	go func() {
		slog.Info(fmt.Sprintf("Server started at %s", app.server.Addr))
		if err := app.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Error starting server: ", "error", err)
			os.Exit(1)
		}
	}()

	<-signalChan
	slog.Info("Shutting down ...")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)

	defer cancel()
	if err := app.server.Shutdown(ctx); err != nil {
		slog.Error("Error shutting down: ", "error", err)
		return err
	}

	return nil
}

func setLogger() {
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
