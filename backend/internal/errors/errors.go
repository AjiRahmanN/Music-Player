package errors

import "fmt"

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
	Op      string `json:"-"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Op != "" {
		return fmt.Sprintf("%s: %s", e.Op, e.Message)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// Error constructors
func NewValidationError(field, message string) *AppError {
	return &AppError{
		Code:    "VALIDATION_ERROR",
		Message: fmt.Sprintf("Validation failed for %s: %s", field, message),
		Details: message,
	}
}

func NewDuplicateError(field, value string) *AppError {
	return &AppError{
		Code:    "DUPLICATE_ERROR",
		Message: fmt.Sprintf("%s '%s' already exists", field, value),
		Details: fmt.Sprintf("The %s '%s' is already registered", field, value),
	}
}

func NewAuthError(message string) *AppError {
	return &AppError{
		Code:    "AUTHENTICATION_ERROR",
		Message: message,
	}
}

func NewNotFoundError(resource, id string) *AppError {
	return &AppError{
		Code:    "NOT_FOUND",
		Message: fmt.Sprintf("%s with ID '%s' not found", resource, id),
	}
}

func NewInternalError(err error) *AppError {
	return &AppError{
		Code:    "INTERNAL_ERROR",
		Message: "An internal server error occurred",
		Err:     err,
	}
}

// Common errors
var (
	ErrInvalidCredentials = NewAuthError("Invalid email or password")
	ErrTokenExpired       = NewAuthError("Token has expired")
	ErrTokenInvalid       = NewAuthError("Invalid token")
)