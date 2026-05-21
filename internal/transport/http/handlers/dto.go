package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type SetProgressRequest struct {
	Slug string `json:"slug"`
}

type ProgressResponse struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProgressResponseList struct {
	ProgressList []ProgressResponse `json:"progress_list"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type CreateProblemRequest struct {
	Name       string `json:"name"`
	Difficulty string `json:"difficulty"`
	Tag        string `json:"tag"`
}

type ProblemResponse struct {
	ID         uuid.UUID  `json:"id"`
	Name       string     `json:"name"`
	Slug       string     `json:"slug"`
	Difficulty string     `json:"difficulty"`
	Tag        string     `json:"tag"`
	Status     string     `json:"status"`
	AuthorID   uuid.UUID  `json:"author_id"`
	VerifiedAt *time.Time `json:"verified_at"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type CreateProblemContentRequest struct {
	ActualGraph   json.RawMessage `json:"actual_graph"`
	ExpectedGraph json.RawMessage `json:"expected_graph"`
	FullText      string          `json:"full_text"`
}

type UpdateProblemContentRequest struct {
	ActualGraph   json.RawMessage `json:"actual_graph"`
	ExpectedGraph json.RawMessage `json:"expected_graph"`
	FullText      string          `json:"full_text"`
}

type SolveProblemRequest struct {
	ActualGraph json.RawMessage `json:"actual_graph"`
}

type SolveProblemResponse struct {
	Solved bool `json:"solved"`
}

type ProblemContentResponse struct {
	ID            uuid.UUID       `json:"id"`
	ProblemID     uuid.UUID       `json:"problem_id"`
	AuthorID      uuid.UUID       `json:"author_id"`
	ActualGraph   json.RawMessage `json:"actual_graph"`
	ExpectedGraph json.RawMessage `json:"expected_graph"`
	FullText      string          `json:"full_text"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}
