package middleware

import (
	"log"
	"net/http"

	"github.com/lmousom/passless-auth/internal/errors"
)

// ErrorHandler is a middleware that handles errors in a consistent way
func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				appErr := errors.NewInternalServer("Internal server error", nil)
				appErr.WriteJSON(w)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// ErrorResponse is a helper function to write error responses
func ErrorResponse(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case *errors.AppError:
		e.WriteJSON(w)
	default:
		appErr := errors.NewInternalServer("Internal server error", err)
		appErr.WriteJSON(w)
	}
}
