package problemcatalog

import (
	"github.com/WebCraftersGH/Education-service/internal/domain"
	"github.com/google/uuid"
	"time"
)

type createProblemRequest struct {
	Name       string `json:"name"`
	Slug       string `json:"slug"`
	Difficulty string `json:"difficulty"`

	Tag      string    `json:"tag"`
	AuthorID uuid.UUID `json:"author_id"`
}

type updateProblemRequest struct {
	Name       string    `json:"name"`
	Slug       string    `json:"slug"`
	Difficulty string    `json:"difficulty"`
	Tag        string    `json:"tag"`
	AuthorID   uuid.UUID `json:"author_id"`
}

type problemResponse struct {
	ID         uuid.UUID            `json:"id"`
	Name       string               `json:"name"`
	Slug       string               `json:"slug"`
	Difficulty string               `json:"difficulty"`
	Tag        string               `json:"tag"`
	Status     domain.ProblemStatus `json:"status"`
	AuthorID   uuid.UUID            `json:"author_id"`
	VerifiedAt *time.Time           `json:"verified_at,omitempty"`
	CreatedAt  time.Time            `json:"created_at"`
	UpdatedAt  time.Time            `json:"updated_at"`
}

type createProblemContentRequest struct {
	DescriptionMD string `json:"description_md"`
	InputFormatMD string `json:"input_format_md"`

	OutputFormatMD string `json:"output_format_md"`
	ConstraintsMD  string `json:"constraints_md"`
	NotesMD        string `json:"notes_md"`
}

type updateProblemContentRequest struct {
	DescriptionMD  string `json:"description_md"`
	InputFormatMD  string `json:"input_format_md"`
	OutputFormatMD string `json:"output_format_md"`
	ConstraintsMD  string `json:"constraints_md"`
	NotesMD        string `json:"notes_md"`
}

type problemContentResponse struct {
	ID            uuid.UUID `json:"id"`
	ProblemID     uuid.UUID `json:"problem_id"`
	DescriptionMD string    `json:"description_md"`

	InputFormatMD  string `json:"input_format_md"`
	OutputFormatMD string `json:"output_format_md"`
	ConstraintsMD  string `json:"constraints_md"`

	NotesMD   string    `json:"notes_md"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
