package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_duration_seconds",
		Help: "Duration of HTTP requests.",
	}, []string{"path", "method", "status"})

	otpRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "otp_requests_total",
		Help: "Total number of OTP requests",
	})

	otpVerifications = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "otp_verifications_total",
		Help: "Total number of OTP verifications",
	}, []string{"status"})
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap response writer to capture status code
		wrapped := wrapResponseWriter(w)
		next.ServeHTTP(wrapped, r)

		duration := time.Since(start).Seconds()
		httpDuration.WithLabelValues(r.URL.Path, r.Method, strconv.Itoa(wrapped.status)).Observe(duration)
	})
}
