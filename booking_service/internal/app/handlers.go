package app

import (
	"net/http"

	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/adapter/booking"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/adapter/middleware"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/adapter/webhook"
	prometheus2 "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/infrastructure/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
)

type Handlers struct {
	BookingHandler *booking.BookingHandler
	AuthHandler    *middleware.AuthMiddleware
	WebhookHandler *webhook.WebhookHandler
	MetricsHandler *middleware.MetricsMiddleware
}

func NewHandlers(services *Services, registry *prometheus.Registry) *Handlers {
	bookingHandler := booking.NewBookingHandler(services.BookingSaver, services.BookingProvider)
	authHandler := middleware.NewAuthMiddleware(services.TokenService)
	webhookHandler := webhook.NewWebhookHandler(services.EventSaver)

	httpMetrics := prometheus2.NewHTTPMetrics(registry)
	metricsHandler := middleware.NewMetricsMiddleware(httpMetrics)

	return &Handlers{bookingHandler, authHandler, webhookHandler, metricsHandler}
}

func (h *Handlers) registerRoutes(cfg *Config, registry *prometheus.Registry) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{Registry: registry}))

	mux.HandleFunc("POST /api/booking", h.BookingHandler.BookRoom)
	mux.HandleFunc("GET /api/booking/hotelier", h.BookingHandler.GetHotelBookings)
	mux.HandleFunc("GET /api/booking/client", h.BookingHandler.GetUserBookings)
	mux.HandleFunc("GET /api/booking/occupied", h.BookingHandler.GetOccupiedRoomDurations)

	mux.HandleFunc("POST "+cfg.HTTPConfig.WebhookHandlerEndpoint, h.WebhookHandler.ServeWebhook)

	handler := h.AuthHandler.Authorize(mux)
	handler = h.MetricsHandler.HandleMetrics(handler)

	c := cors.New(cors.Options{
		AllowedOrigins:   cfg.HTTPConfig.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	return c.Handler(handler)
}
