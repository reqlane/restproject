package apperrors

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrMissingID    = errors.New("missing id")
	ErrInvalidID    = errors.New("invalid id")
	ErrInvalidField = errors.New("invalid field type")
)
