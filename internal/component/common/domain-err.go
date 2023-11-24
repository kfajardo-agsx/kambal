package common

import (
	"errors"
	"strings"

	"github.com/jinzhu/gorm"
)

type (
	APIError struct {
		Type ErrorType
		Msg  string
	}
)

func (e APIError) Error() string {
	return e.Msg
}

// ErrorType describes the type of error
type ErrorType int

const (
	ErrorTypeUnreachable ErrorType = iota
	ErrorTypeNotFound
	ErrorTypeConflict
	ErrorTypeBadRequest
	ErrorTypeUnauthorized
	ErrorTypeForbidden
	ErrorTypeUnknown
	ErrorTypeServiceUnavailable
	ErrorTypeNotAllowed
	ErrorExternalService
)

func NewAPIError(t ErrorType, msg string) error {
	return &APIError{t, msg}
}

func RepositoryErrorToAPIError(err error) error {
	msg := err.Error()
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return NewAPIError(ErrorTypeNotFound, msg)
	case strings.Contains(msg, "already exists"):
		return NewAPIError(ErrorTypeConflict, msg)
	case strings.Contains(msg, "not allowed"):
		return NewAPIError(ErrorTypeNotAllowed, msg)
	case strings.Contains(msg, "does not exist"):
		return NewAPIError(ErrorTypeBadRequest, msg)
	default:
		return NewAPIError(ErrorTypeUnknown, msg)
	}
}
