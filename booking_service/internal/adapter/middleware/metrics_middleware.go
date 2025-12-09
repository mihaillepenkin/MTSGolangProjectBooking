package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	prometheus2 "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/infrastructure/prometheus"
)

var (
	endpoints = []string{
		"/api/booking",
		"/api/booking/hotelier",
		"/api/booking/client",
		"/api/booking/occupied",
	}
)

type MetricsMiddleware struct {
	metrics *prometheus2.HTTPMetrics
}

func NewMetricsMiddleware(metrics *prometheus2.HTTPMetrics) *MetricsMiddleware {
	return &MetricsMiddleware{metrics: metrics}
}

func (m *MetricsMiddleware) HandleMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rw, r)

		status := rw.status

		statusGroup := getStatusGroup(status)

		endpoint, err := getEndpoint(r.URL.Path)
		if err != nil {
			return
		}

		m.metrics.RequestsTotal.WithLabelValues(r.Method, endpoint, statusGroup).Inc()
		m.metrics.RequestsDuration.WithLabelValues(r.Method, endpoint, statusGroup).Observe(time.Since(start).Seconds())
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (w *responseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func getStatusGroup(code int) string {
	switch {
	case code >= 200 && code <= 299:
		return "2xx"
	case code >= 400 && code <= 499:
		return "4xx"
	case code >= 500 && code <= 599:
		return "5xx"
	default:
		return "unknown"
	}
}

func getEndpoint(endpoint string) (string, error) {
	maxCommonUrl := ""
	for _, ep := range endpoints {
		if len(endpoint) >= len(ep) && strings.HasPrefix(endpoint, ep) && len(maxCommonUrl) < len(endpoint) {
			maxCommonUrl = ep
		}
	}
	if len(maxCommonUrl) == 0 {
		return "", errors.New("endpoint not found")
	}
	return maxCommonUrl, nil
}
