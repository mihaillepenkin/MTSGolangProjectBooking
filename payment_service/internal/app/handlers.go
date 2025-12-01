package app

import (
	"net/http"

	"github.com/mihaillepenkin/MTSGolangProjectBooking/payment_service/internal/adapter"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/payment_service/internal/adapter/httpconfig"
	"github.com/rs/cors"
)

type Handlers struct {
	PaymentHandler *adapter.PaymentHandler
}

func NewHandlers(cfg httpconfig.HTTPConfig, services *Services) *Handlers {
	paymentHandler := adapter.NewPaymentHandler(services.PaymentService, cfg)
	return &Handlers{PaymentHandler: paymentHandler}
}

func (h *Handlers) registerRoutes(cfg *Config) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/payment", h.PaymentHandler.CreatePayment)
	mux.HandleFunc("POST /api/payment/process", h.PaymentHandler.ProcessPayment)

	c := cors.New(cors.Options{
		AllowedOrigins:   cfg.HTTPConfig.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	return c.Handler(mux)
}
