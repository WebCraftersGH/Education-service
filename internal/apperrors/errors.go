package apperrors

import (
	"errors"
)

var (
	ErrDuplicateRecord = errors.New("repository duplicate error")	
	ErrInternal = errors.New("repository internal error")
	ErrInvalidArgument = errors.New("invalid argument")
)

