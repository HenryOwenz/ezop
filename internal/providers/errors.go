package providers

import (
	"fmt"
)

// Common errors
var (
	ErrProviderNotFound  = fmt.Errorf("provider not found")
	ErrNotAuthenticated  = fmt.Errorf("provider not authenticated")
	ErrInvalidConfig     = fmt.Errorf("invalid configuration")
	ErrOperationNotFound = fmt.Errorf("operation not found")
	ErrServiceNotFound   = fmt.Errorf("service not found")
	ErrCategoryNotFound  = fmt.Errorf("category not found")
	ErrNotImplemented    = fmt.Errorf("not implemented")
)

// ProviderError represents an error from a provider
type ProviderError struct {
	Provider string
	Err      error
}

// Error returns the error message
func (e *ProviderError) Error() string {
	return fmt.Sprintf("%s: %s", e.Provider, e.Err)
}

// Unwrap returns the underlying error
func (e *ProviderError) Unwrap() error {
	return e.Err
}

// NewProviderError creates a new provider error
func NewProviderError(provider string, err error) error {
	return &ProviderError{
		Provider: provider,
		Err:      err,
	}
}

// OperationError represents an error from an operation
type OperationError struct {
	Operation string
	Err       error
}

// Error returns the error message
func (e *OperationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Operation, e.Err)
}

// Unwrap returns the underlying error
func (e *OperationError) Unwrap() error {
	return e.Err
}

// NewOperationError creates a new operation error
func NewOperationError(operation string, err error) error {
	return &OperationError{
		Operation: operation,
		Err:       err,
	}
}
