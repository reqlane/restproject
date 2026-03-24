package apperrors

import (
	"errors"
	"net/http"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrMissingID    = errors.New("missing id")
	ErrInvalidID    = errors.New("invalid id")
	ErrInvalidField = errors.New("invalid field type")
)

type Error struct {
	domainError error
	appError    error
}

func (e *Error) DomainError() error {
	return e.domainError
}

func (e *Error) AppError() error {
	return e.appError
}

func (e *Error) Error() string {
	return e.appError.Error()
}

func NewError(domainError, appError error) error {
	return &Error{
		domainError: domainError,
		appError:    appError,
	}
}

type HTTPError struct {
	Status  int
	Message string
}

func FromError(err error) *HTTPError {
	if err, ok := errors.AsType[*Error](err); ok {
		httpError := &HTTPError{Message: err.AppError().Error()}
		switch err.DomainError() {
		case ErrNotFound:
			httpError.Status = http.StatusNotFound
		case ErrMissingID, ErrInvalidID, ErrInvalidField:
			httpError.Status = http.StatusBadRequest
		default:
			httpError.Status = http.StatusInternalServerError
			httpError.Message = "Internal server error"
		}
		return httpError
	}
	return &HTTPError{
		Status:  http.StatusInternalServerError,
		Message: "Internal server error",
	}
}
