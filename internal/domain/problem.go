package domain

import (
	"github.com/google/uuid"
	"time"
)

type ProblemStatus string

const (
	ProblemStatusDraft    ProblemStatus = "draft"
	ProblemStatusApproved ProblemStatus = "approved"
	ProblemStatusRejected ProblemStatus = "rejected"
)

func (s ProblemStatus) IsValid() bool {
	switch s {
	case ProblemStatusDraft, ProblemStatusApproved, ProblemStatusRejected:
		return true
	default:
		return false
	}
}

type Problem struct {
	ID         uuid.UUID
	Name       string
	Slug       string
	Difficulty string
	Tag        string
	Status     ProblemStatus
	AuthorID   uuid.UUID
	VerifiedAt *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type ProblemContent struct {
	ID             uuid.UUID
	ProblemID      uuid.UUID
	DescriptionMD  string
	InputFormatMD  string
	OutputFormatMD string
	ConstraintsMD  string
	NotesMD        string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type ProblemFilter struct {
	Tag        string
	Difficulty string
	Limit      int
	Offset     int
}
