package base

import (
	"fmt"
	"net/http"
)

// ServiceError is a general optional error that can be
// returned by any type of service. NOT SERIALIZABLE.
type ServiceError struct {
	Message string
	Blame   Blame
	Code    int
	Err     error
}

// NewPostgresWriteError returns ServiceError with general write error message.
func NewPostgresWriteError(err error) *ServiceError {
	return &ServiceError{
		Err:     err,
		Blame:   BlamePostgres,
		Code:    http.StatusInternalServerError,
		Message: "failed to write data to database",
	}
}

// NewPostgresReadError returns ServiceError with general read error message.
func NewPostgresReadError(err error) *ServiceError {
	return &ServiceError{
		Err:     err,
		Blame:   BlamePostgres,
		Code:    http.StatusInternalServerError,
		Message: "failed to read data from database",
	}
}

// NewNeo4jReadError returns ServiceError with general read error message.
func NewNeo4jReadError(err error) *ServiceError {
	return &ServiceError{
		Err:     err,
		Blame:   BlameNeo4j,
		Code:    http.StatusInternalServerError,
		Message: "failed to read data from neo4j",
	}
}

// NewNeo4jWriteError returns ServiceError with general write error message.
func NewNeo4jWriteError(err error) *ServiceError {
	return &ServiceError{
		Err:     err,
		Blame:   BlameNeo4j,
		Code:    http.StatusInternalServerError,
		Message: "failed to read data from neo4j",
	}
}

// NewNotFoundError returns ServiceError with general not found error message.
func NewNotFoundError(err error) *ServiceError {
	return &ServiceError{
		Err:     err,
		Blame:   BlameUser,
		Code:    http.StatusNotFound,
		Message: "not found",
	}
}

func NewParseEnumError(err error) *ServiceError {
	return &ServiceError{
		Err:     err,
		Blame:   BlameUser,
		Code:    http.StatusInternalServerError,
		Message: "failed to parse enum",
	}
}

func NewParseJWTError(err error) *ServiceError {
	return &ServiceError{
		Err:     err,
		Blame:   BlameServer,
		Code:    http.StatusInternalServerError,
		Message: "failed to parse jwt-token",
	}
}

func NewGenerateJWTError(err error) *ServiceError {
	return &ServiceError{
		Err:     err,
		Blame:   BlameServer,
		Code:    http.StatusInternalServerError,
		Message: "failed to generate jwt-token",
	}
}

func (e *ServiceError) Error() string {
	return fmt.Sprintf("[%d] %v (blame: %s)", e.Code, e.Err, e.Blame)
}
