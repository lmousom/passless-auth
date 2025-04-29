package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ErrorCode represents a specific type of error
type ErrorCode string

const (
	// Common error codes
	ErrInvalidRequest     ErrorCode = "INVALID_REQUEST"
	ErrUnauthorized       ErrorCode = "UNAUTHORIZED"
	ErrForbidden          ErrorCode = "FORBIDDEN"
	ErrNotFound           ErrorCode = "NOT_FOUND"
	ErrInternalServer     ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE"

	// Auth specific error codes
	ErrInvalidOTP      ErrorCode = "INVALID_OTP"
	ErrOTPExpired      ErrorCode = "OTP_EXPIRED"
	ErrTooManyAttempts ErrorCode = "TOO_MANY_ATTEMPTS"
	ErrInvalidToken    ErrorCode = "INVALID_TOKEN"
	ErrTokenExpired    ErrorCode = "TOKEN_EXPIRED"
)

// AppError represents an application error
type AppError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Err     error     `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// New creates a new AppError
func New(code ErrorCode, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// HTTPStatus returns the appropriate HTTP status code for the error
func (e *AppError) HTTPStatus() int {
	switch e.Code {
	case ErrInvalidRequest:
		return http.StatusBadRequest
	case ErrUnauthorized:
		return http.StatusUnauthorized
	case ErrForbidden:
		return http.StatusForbidden
	case ErrNotFound:
		return http.StatusNotFound
	case ErrInternalServer:
		return http.StatusInternalServerError
	case ErrServiceUnavailable:
		return http.StatusServiceUnavailable
	case ErrInvalidOTP, ErrOTPExpired, ErrTooManyAttempts, ErrInvalidToken, ErrTokenExpired:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

// WriteJSON writes the error as JSON to the response writer
func (e *AppError) WriteJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.HTTPStatus())

	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]interface{}{
			"code":    e.Code,
			"message": e.Message,
		},
	})
}

// Helper functions for common errors
func NewInvalidRequest(message string, err error) *AppError {
	return New(ErrInvalidRequest, message, err)
}

func NewUnauthorized(message string, err error) *AppError {
	return New(ErrUnauthorized, message, err)
}

func NewForbidden(message string, err error) *AppError {
	return New(ErrForbidden, message, err)
}

func NewNotFound(message string, err error) *AppError {
	return New(ErrNotFound, message, err)
}

func NewInternalServer(message string, err error) *AppError {
	return New(ErrInternalServer, message, err)
}

func NewServiceUnavailable(message string, err error) *AppError {
	return New(ErrServiceUnavailable, message, err)
}

func NewInvalidOTP(message string, err error) *AppError {
	return New(ErrInvalidOTP, message, err)
}

func NewOTPExpired(message string, err error) *AppError {
	return New(ErrOTPExpired, message, err)
}

func NewTooManyAttempts(message string, err error) *AppError {
	return New(ErrTooManyAttempts, message, err)
}

func NewInvalidToken(message string, err error) *AppError {
	return New(ErrInvalidToken, message, err)
}

func NewTokenExpired(message string, err error) *AppError {
	return New(ErrTokenExpired, message, err)
}
