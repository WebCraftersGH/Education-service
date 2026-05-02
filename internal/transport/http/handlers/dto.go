package handlers

import (
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
