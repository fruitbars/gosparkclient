package gosparkclient

import (
	"fmt"
)

// ErrorType represents the type of error that occurred
type ErrorType string

const (
	ErrConfiguration  ErrorType = "ConfigurationError"
	ErrConnection     ErrorType = "ConnectionError"
	ErrAuthentication ErrorType = "AuthenticationError"
	ErrRequest        ErrorType = "RequestError"
	ErrResponse       ErrorType = "ResponseError"
	ErrWebSocket      ErrorType = "WebSocketError"
)

// SparkError represents a custom error type for the Spark client
type SparkError struct {
	Type    ErrorType
	Message string
	Err     error
}

// Error implements the error interface
func (e *SparkError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (underlying: %v)", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Unwrap returns the underlying error
func (e *SparkError) Unwrap() error {
	return e.Err
}

// NewSparkError creates a new SparkError
func NewSparkError(errType ErrorType, message string, err error) *SparkError {
	return &SparkError{
		Type:    errType,
		Message: message,
		Err:     err,
	}
}

// Helper functions for creating specific error types
func newConfigError(message string, err error) *SparkError {
	return NewSparkError(ErrConfiguration, message, err)
}

func newConnectionError(message string, err error) *SparkError {
	return NewSparkError(ErrConnection, message, err)
}

func newAuthError(message string, err error) *SparkError {
	return NewSparkError(ErrAuthentication, message, err)
}

func newRequestError(message string, err error) *SparkError {
	return NewSparkError(ErrRequest, message, err)
}

func newResponseError(message string, err error) *SparkError {
	return NewSparkError(ErrResponse, message, err)
}

func newWebSocketError(message string, err error) *SparkError {
	return NewSparkError(ErrWebSocket, message, err)
}
