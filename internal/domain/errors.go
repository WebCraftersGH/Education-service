package domain

import (
	"errors"
)

var (
	ErrProblemNameRequired             = errors.New("problem name is required")
	ErrProblemSlugRequired             = errors.New("problem slug is required")
	ErrProblemDifficultyRequired       = errors.New("problem difficulty is required")
	ErrProblemAuthorIDRequired         = errors.New("problem author id is required")
	ErrProblemIDRequired               = errors.New("problem id is required")
	ErrProblemContentIDRequired        = errors.New("problem content id is required")
	ErrProblemContentProblemIDRequired = errors.New("problem id is required")
	ErrProblemDescriptionRequired      = errors.New("problem description is required")
)
