package main

import (
	"github.com/mihaillepenkin/MTSGolangProjectBooking/payment_service/internal/app"
)

func main() {
	application, err := app.NewApp()
	if err != nil {
		return
	}

	err = application.Start()
	if err != nil {
		return
	}
}
