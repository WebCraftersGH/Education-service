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

func toProgressResponseList(checkpoints []domain.CheckPoint) ProgressResponseList {
	list := make([]ProgressResponse, len(checkpoints))
	for i, checkpoint := range checkpoints {
		ck := toProgressResponse(checkpoint)
		list[i] = ck
	}
	return ProgressResponseList{ProgressList: list}
}
