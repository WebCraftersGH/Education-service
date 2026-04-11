package courseprogress

import (
	"errors"
)

var (
	ErrDuplicateRecord = errors.New("repository duplicate error")	
	ErrInternal = errors.New("repository internal error")
)
