package main

import (
	"log/slog"

	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/app"
)

func main() {
	application, err := app.NewApp()
	if err != nil {
		slog.Error("Error creating new application", "error", err)
		return
	}

	err = application.Start()
	if err != nil {
		slog.Error("Error starting application", "error", err)
	}
}
