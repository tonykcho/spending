package middlewares

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests",
		},
		[]string{"method", "path"},
	)
)

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Wrap ResponseWriter to capture status code
		next.ServeHTTP(w, r)

		httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path).Inc()
		// Optionally record duration as a histogram
	})
}
