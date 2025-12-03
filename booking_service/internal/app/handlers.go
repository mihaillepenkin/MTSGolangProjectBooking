package app

import (
	"net/http"

	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/adapter/booking"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/adapter/middleware"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/adapter/webhook"
	"github.com/rs/cors"
)

type Handlers struct {
	BookingHandler *booking.BookingHandler
	AuthHandler    *middleware.AuthMiddleware
	WebhookHandler *webhook.WebhookHandler
}

func NewHandlers(services *Services) *Handlers {
	bookingHandler := booking.NewBookingHandler(services.BookingSaver, services.BookingProvider)
	authHandler := middleware.NewAuthMiddleware(services.TokenService)
	webhookHandler := webhook.NewWebhookHandler(services.BookingSaver)
	return &Handlers{bookingHandler, authHandler, webhookHandler}
}

func (h *Handlers) registerRoutes(cfg *Config) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/booking", h.BookingHandler.BookRoom)
	mux.HandleFunc("GET /api/booking/hotelier", h.BookingHandler.GetHotelBookings)
	mux.HandleFunc("GET /api/booking/client", h.BookingHandler.GetUserBookings)
	mux.HandleFunc("GET /api/booking/occupied", h.BookingHandler.GetOccupiedRoomDurations)

	mux.HandleFunc("POST "+cfg.HTTPConfig.WebhookHandlerEndpoint, h.WebhookHandler.ServeWebhook)

	handler := h.AuthHandler.Authorize(mux)

	c := cors.New(cors.Options{
		AllowedOrigins:   cfg.HTTPConfig.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	return c.Handler(handler)
}
