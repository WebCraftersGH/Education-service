package domain

import (
	"time"

	"github.com/google/uuid"
)

type CheckPoint struct {
	ID uuid.UUID
	UserID uuid.UUID
	Slug string
	CreatedAt time.Time
	UpdatedAt time.Time
}
