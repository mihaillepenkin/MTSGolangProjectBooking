package prometheus

import "github.com/prometheus/client_golang/prometheus"

type HTTPMetrics struct {
	RequestsTotal    *prometheus.CounterVec
	RequestsDuration *prometheus.HistogramVec
}

func NewHTTPMetrics(reg prometheus.Registerer) *HTTPMetrics {
	m := &HTTPMetrics{
		RequestsTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number og HTTP requests",
		}, []string{"method", "endpoint", "status"}),
		RequestsDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "endpoint", "status"}),
	}

	reg.MustRegister(m.RequestsTotal)
	reg.MustRegister(m.RequestsDuration)
	return m
}
