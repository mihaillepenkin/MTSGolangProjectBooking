package main

import (
	"context"
	"hotel_service/internal/infrastructure/api"
	"hotel_service/internal/infrastructure/config"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	db := config.ConfigureDb()
	if db == nil {
		slog.Error("Error connecting to database")
	}
	err := db.Ping()
	if err != nil {
		slog.Error(err.Error())
	}
	slog.Info("Successfully connected to database")

	hotelHandler := api.HotelHandler{}
	hotelHandler.Initialize(db)

	mux := api.CreateRouting(&hotelHandler)
	handler := api.CORSMiddleware(api.AuthMiddleware(mux))

	server := &http.Server{
		Addr:    ":8082",
		Handler: handler,
	}
	go func() {
		err := server.ListenAndServe()
		slog.Info("Server is running on port 8082...")
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
