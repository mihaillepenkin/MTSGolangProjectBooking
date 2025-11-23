package main

import (
	"context"
	"database/sql"
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
	db, err := config.ConfigureDb()
	if err != nil {
		slog.Error(err.Error())
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			slog.Error(err.Error())
		}
	}(db)
	err = db.Ping()
	if err != nil {
		slog.Error(err.Error())
	}

	hotelHandler := api.HotelHandler{}
	hotelHandler.Initialize(db)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1/hotels", hotelHandler.GetAllHotels)

	mux.HandleFunc("POST /api/v1/hotels", hotelHandler.AddHotelInfo)

	mux.HandleFunc("PUT /api/v1/hotels", hotelHandler.UpdateHotelInfo)

	mux.HandleFunc("GET /api/v1/hotels/{id}", hotelHandler.GetHotelById)

	server := &http.Server{
		Addr:    ":8082",
		Handler: mux,
	}
	go func() {
		err := server.ListenAndServe()
		slog.Debug("Server is running on port 8082...")
		if err != nil {
			slog.Error(err.Error())
		}
	}()

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)
	<-osSignals
	slog.Debug("Graceful shutdown...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = server.Shutdown(ctx)
	if err != nil {
		slog.Error(err.Error())
	}
}
