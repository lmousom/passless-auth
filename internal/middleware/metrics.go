package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP metrics
	httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_duration_seconds",
		Help: "Duration of HTTP requests.",
	}, []string{"path", "method", "status"})

	httpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	}, []string{"path", "method", "status"})

	httpRequestsInFlight = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "http_requests_in_flight",
		Help: "Current number of HTTP requests being served",
	})

	// OTP metrics
	otpRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "otp_requests_total",
		Help: "Total number of OTP requests",
	}, []string{"status"})

	otpVerifications = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "otp_verifications_total",
		Help: "Total number of OTP verifications",
	}, []string{"status"})

	// 2FA metrics
	twoFAEnrollments = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "twofa_enrollments_total",
		Help: "Total number of 2FA enrollments",
	}, []string{"status"})

	twoFAVerifications = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "twofa_verifications_total",
		Help: "Total number of 2FA verifications",
	}, []string{"status"})

	// Redis metrics
	redisOperations = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "redis_operations_total",
		Help: "Total number of Redis operations",
	}, []string{"operation", "status"})

	redisOperationDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "redis_operation_duration_seconds",
		Help: "Duration of Redis operations",
	}, []string{"operation"})
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
		// Track in-flight requests
		httpRequestsInFlight.Inc()
		defer httpRequestsInFlight.Dec()

		start := time.Now()

		// Wrap response writer to capture status code
		wrapped := wrapResponseWriter(w)
		next.ServeHTTP(wrapped, r)

		// Record metrics
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(wrapped.status)

		// Record HTTP metrics
		httpDuration.WithLabelValues(r.URL.Path, r.Method, status).Observe(duration)
		httpRequestsTotal.WithLabelValues(r.URL.Path, r.Method, status).Inc()

		// Record specific metrics based on the endpoint
		switch r.URL.Path {
		case "/api/v1/sendOtp":
			otpRequests.WithLabelValues(status).Inc()
		case "/api/v1/verifyOtp":
			otpVerifications.WithLabelValues(status).Inc()
		case "/api/v1/2fa/enable":
			twoFAEnrollments.WithLabelValues(status).Inc()
		case "/api/v1/2fa/verify":
			twoFAVerifications.WithLabelValues(status).Inc()
		}
	})
}
