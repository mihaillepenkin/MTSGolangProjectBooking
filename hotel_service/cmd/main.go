package main

import (
	"context"
	"hotel_service/internal/config"
	"hotel_service/internal/infrastructure/grpc"
	"hotel_service/internal/infrastructure/http"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	textHandler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger := slog.New(textHandler)
	slog.SetDefault(logger)

	cfg := config.LoadConfig()
	if cfg == nil {
		slog.Error("Error loading config")
		return
	}
	db := cfg.ConnectToDb()
	if db == nil {
		slog.Error("Error connecting to database")
		return
	}
	err := cfg.RunMigrations(db)
	if err != nil {
		slog.Error("Error running migrations")
		return
	}

	hotelGrpcService := grpc2.HotelGrpcService{}
	hotelGrpcService.Initialize(db)

	grpcServer, listener := cfg.ConfigureGrpcServer(&hotelGrpcService)
	if grpcServer == nil || listener == nil {
		slog.Error("Error configuring gRPC server")
		return
	}
	err = grpcServer.Serve(listener)
	if err != nil {
		slog.Error("Error running gRPC server")
		return
	}
	slog.Info("gRPC server is running...")

	hotelHandler := api.HotelHandler{}
	hotelHandler.Initialize(db)

	mux := api.CreateRouting(&hotelHandler)
	handler := api.AuthMiddleware(api.CORSMiddleware(mux))

	server := cfg.ConfigureHttpServer(&handler)
	go func() {
		err := server.ListenAndServe()
		slog.Info("Http server is running...")
		if err != nil {
			slog.Error(err.Error())
		}
	}()

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)
	<-osSignals
	slog.Info("Graceful shutdown...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = server.Shutdown(ctx)
	if err != nil {
		slog.Error(err.Error())
	}
}
