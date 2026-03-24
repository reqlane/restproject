package apperrors

import "errors"

var (
	ErrNotFound     = errors.New("not found") // TODO custom error for specified id
	ErrMissingID    = errors.New("missing id")
	ErrInvalidID    = errors.New("invalid id")
	ErrInvalidField = errors.New("invalid field type") // TODO custom error for specified field name
)
