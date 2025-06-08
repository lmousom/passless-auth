package middleware

import (
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func TracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tracer := otel.Tracer("http")

		// Extract request details for attributes
		attrs := []attribute.KeyValue{
			attribute.String("http.method", r.Method),
			attribute.String("http.url", r.URL.String()),
			attribute.String("http.user_agent", r.UserAgent()),
			attribute.String("http.remote_addr", r.RemoteAddr),
			attribute.String("http.host", r.Host),
		}

		// Add query parameters as attributes
		for key, values := range r.URL.Query() {
			if len(values) > 0 {
				attrs = append(attrs, attribute.String("http.query."+key, values[0]))
			}
		}

		// Start span with attributes
		ctx, span := tracer.Start(ctx, r.URL.Path,
			trace.WithAttributes(attrs...),
			trace.WithSpanKind(trace.SpanKindServer),
		)
		defer span.End()

		// Create a custom response writer to capture status code
		rw := &tracingResponseWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		// Add start time for duration calculation
		start := time.Now()

		// Process request
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)

		// Record additional span attributes after request processing
		span.SetAttributes(
			attribute.Int("http.status_code", rw.status),
			attribute.Float64("http.duration_ms", float64(time.Since(start).Milliseconds())),
		)

		// Set span status based on HTTP status code
		if rw.status >= 400 {
			span.SetStatus(codes.Error, http.StatusText(rw.status))
		} else {
			span.SetStatus(codes.Ok, "")
		}
	})
}

// tracingResponseWriter is a custom response writer that captures the status code
type tracingResponseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *tracingResponseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
