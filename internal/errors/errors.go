package errors

import (
	"errors"
	"fmt"
)

// Error types for better error handling
var (
	// Validation errors
	ErrInvalidInput      = errors.New("invalid input")
	ErrEmptyField        = errors.New("required field is empty")
	ErrInvalidFormat     = errors.New("invalid format")
	ErrNegativeValue     = errors.New("value cannot be negative")
	
	// Resource errors
	ErrNotFound          = errors.New("resource not found")
	ErrAlreadyExists     = errors.New("resource already exists")
	
	// Database errors
	ErrDatabaseOperation = errors.New("database operation failed")
	ErrMigration         = errors.New("migration failed")
	ErrConnection        = errors.New("connection failed")
)

// ValidationError represents a validation error with field information
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string) error {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

// NotFoundError represents a resource not found error
type NotFoundError struct {
	Resource string
	ID       string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with ID '%s' not found", e.Resource, e.ID)
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(resource, id string) error {
	return &NotFoundError{
		Resource: resource,
		ID:       id,
	}
}

// DatabaseError represents a database operation error
type DatabaseError struct {
	Operation string
	Err       error
}

func (e *DatabaseError) Error() string {
	return fmt.Sprintf("database %s failed: %v", e.Operation, e.Err)
}

func (e *DatabaseError) Unwrap() error {
	return e.Err
}

// NewDatabaseError creates a new database error
func NewDatabaseError(operation string, err error) error {
	return &DatabaseError{
		Operation: operation,
		Err:       err,
	}
}

// IsValidationError checks if an error is a validation error
func IsValidationError(err error) bool {
	var validationErr *ValidationError
	return errors.As(err, &validationErr)
}

// IsNotFoundError checks if an error is a not found error
func IsNotFoundError(err error) bool {
	var notFoundErr *NotFoundError
	return errors.As(err, &notFoundErr)
}

// IsDatabaseError checks if an error is a database error
func IsDatabaseError(err error) bool {
	var dbErr *DatabaseError
	return errors.As(err, &dbErr)
}

