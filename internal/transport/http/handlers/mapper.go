package handlers

import (
	"github.com/WebCraftersGH/Education-service/internal/domain"
)

func toProgressResponse(checkpoint domain.CheckPoint) ProgressResponse {
	return ProgressResponse{
		ID:        checkpoint.ID,
		UserID:    checkpoint.UserID,
		Slug:      checkpoint.Slug,
		CreatedAt: checkpoint.CreatedAt,
		UpdatedAt: checkpoint.UpdatedAt,
	}
}
