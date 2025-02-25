package errors

import (
	"net/http"
)

// CustomError defines a standardized error response
type CustomError struct {
	Code    string `json:"error_code"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// Predefined errors
var (
	ErrUserNotFound = CustomError{"AUTH_USER_NOT_FOUND", "User not found", http.StatusNotFound}
	ErrInvalidToken = CustomError{"AUTH_INVALID_TOKEN", "Invalid authentication token", http.StatusUnauthorized}
	ErrUnauthorized = CustomError{"AUTH_UNAUTHORIZED", "Unauthorized access", http.StatusForbidden}
	ErrInvalidInput = CustomError{"INPUT_INVALID", "Invalid input provided", http.StatusBadRequest}
	ErrClientNotFound = CustomError{"BUSINESS_CLIENT_NOT_FOUND", "Client profile not found", http.StatusNotFound}
	ErrServerError = CustomError{"SERVER_ERROR", "Internal server error", http.StatusInternalServerError}
)

// ErrorMap allows lookup by error message
var ErrorMap = map[string]CustomError{
	"UserNotFound":      ErrUserNotFound,
	"InvalidToken":      ErrInvalidToken,
	"Unauthorized":      ErrUnauthorized,
	"InvalidInput":      ErrInvalidInput,
	"ClientNotFound":    ErrClientNotFound,
	"InternalError":     ErrServerError,
}

// GetError returns an error struct based on a string key
func GetError(key string) CustomError {
	if err, exists := ErrorMap[key]; exists {
		return err
	}
	return ErrServerError // Default to internal server error
}
