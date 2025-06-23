package tvwsclient

import (
	"errors"
	"fmt"
)

// Error types for better error handling
var (
	ErrConnectionClosed     = errors.New("websocket connection is closed")
	ErrAuthenticationFailed = errors.New("authentication failed")
	ErrInvalidMessage       = errors.New("invalid message format")
	ErrSessionNotFound      = errors.New("session not found")
	ErrInvalidSymbol        = errors.New("invalid symbol format")
	ErrRateLimitExceeded    = errors.New("rate limit exceeded")
	ErrReconnectFailed      = errors.New("reconnection failed")
	ErrTimeout              = errors.New("operation timeout")
)

// TradingViewError wraps errors with additional context
type TradingViewError struct {
	Op      string // operation that failed
	Code    string // error code
	Message string // human-readable message
	Err     error  // underlying error
}

func (e *TradingViewError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%s)", e.Op, e.Message, e.Err.Error())
	}
	return fmt.Sprintf("%s: %s", e.Op, e.Message)
}

func (e *TradingViewError) Unwrap() error {
	return e.Err
}

// Error codes
const (
	ErrCodeConnection    = "CONNECTION_ERROR"
	ErrCodeAuth          = "AUTH_ERROR"
	ErrCodeMessage       = "MESSAGE_ERROR"
	ErrCodeSession       = "SESSION_ERROR"
	ErrCodeSymbol        = "SYMBOL_ERROR"
	ErrCodeRateLimit     = "RATE_LIMIT_ERROR"
	ErrCodeReconnect     = "RECONNECT_ERROR"
	ErrCodeTimeout       = "TIMEOUT_ERROR"
	ErrCodeInternal      = "INTERNAL_ERROR"
	ErrCodeValidation    = "VALIDATION_ERROR"
)

// NewTradingViewError creates a new TradingViewError
func NewTradingViewError(op, code, message string, err error) *TradingViewError {
	return &TradingViewError{
		Op:      op,
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// WrapConnectionError wraps connection-related errors
func WrapConnectionError(op string, err error) error {
	return NewTradingViewError(op, ErrCodeConnection, "connection error", err)
}

// WrapAuthError wraps authentication-related errors
func WrapAuthError(op string, err error) error {
	return NewTradingViewError(op, ErrCodeAuth, "authentication error", err)
}

// WrapMessageError wraps message parsing errors
func WrapMessageError(op string, err error) error {
	return NewTradingViewError(op, ErrCodeMessage, "message processing error", err)
}

// WrapSessionError wraps session-related errors
func WrapSessionError(op string, err error) error {
	return NewTradingViewError(op, ErrCodeSession, "session error", err)
}

// WrapValidationError wraps validation errors
func WrapValidationError(op, message string, err error) error {
	return NewTradingViewError(op, ErrCodeValidation, message, err)
}

// IsConnectionError checks if error is connection-related
func IsConnectionError(err error) bool {
	var tvErr *TradingViewError
	if errors.As(err, &tvErr) {
		return tvErr.Code == ErrCodeConnection
	}
	return errors.Is(err, ErrConnectionClosed)
}

// IsAuthError checks if error is authentication-related
func IsAuthError(err error) bool {
	var tvErr *TradingViewError
	if errors.As(err, &tvErr) {
		return tvErr.Code == ErrCodeAuth
	}
	return errors.Is(err, ErrAuthenticationFailed)
}

// IsSessionError checks if error is session-related
func IsSessionError(err error) bool {
	var tvErr *TradingViewError
	if errors.As(err, &tvErr) {
		return tvErr.Code == ErrCodeSession
	}
	return errors.Is(err, ErrSessionNotFound)
}

// IsRetryableError checks if an error is retryable
func IsRetryableError(err error) bool {
	var tvErr *TradingViewError
	if errors.As(err, &tvErr) {
		switch tvErr.Code {
		case ErrCodeConnection, ErrCodeTimeout, ErrCodeRateLimit:
			return true
		}
	}
	
	// Check for specific error types
	return errors.Is(err, ErrConnectionClosed) || 
		   errors.Is(err, ErrTimeout) || 
		   errors.Is(err, ErrRateLimitExceeded)
}