package courseprogress

import (
	"context"
	"fmt"
	"strings"

	"github.com/WebCraftersGH/Education-service/internal/apperrors"
	"github.com/WebCraftersGH/Education-service/internal/contracts"
	"github.com/WebCraftersGH/Education-service/internal/domain"
	"github.com/google/uuid"
)

const (
	MAX_PROGRESS_LIMIT = 100 //TODO Узкое горлышко
	BASE_GET_PROGRESS_LIMIT = 10
)

type courseProgress struct {
	repo contracts.ProgressRepo
}

func NewCourseProgress(
	repo contracts.ProgressRepo,
) *courseProgress {
	return &courseProgress{
		repo: repo,
	}
}

func (c *courseProgress) SetProgress(
	ctx context.Context, 
	userID uuid.UUID,
	slug string,	
) (domain.CheckPoint, error) {

	slug = strings.TrimSpace(slug)

	if slug == "" {
		return domain.CheckPoint{}, fmt.Errorf("set progress: %w", apperrors.ErrInvalidArgument)
	}

	if userID == uuid.Nil {
		return domain.CheckPoint{}, fmt.Errorf("set progress: %w", apperrors.ErrInvalidArgument)
	}

	checkPoint := domain.CheckPoint{
		UserID: userID,
		Slug: slug,
	}

	newCheckPoint, err := c.repo.CreateCheckPoint(ctx, checkPoint)
	if err != nil {
		return domain.CheckPoint{}, fmt.Errorf("set progress: %w", err)
	}

	return newCheckPoint, nil
}

func (c *courseProgress) ProgressList(
	ctx context.Context,
	userID uuid.UUID,
	limit,
	offset int,
) ([]domain.CheckPoint, error) {
	
	if userID == uuid.Nil {
		return nil, fmt.Errorf("get progress list: %w", apperrors.ErrInvalidArgument)
	}

	if limit <= 0 {
		limit = BASE_GET_PROGRESS_LIMIT
	}

	if limit > MAX_PROGRESS_LIMIT {
		limit = MAX_PROGRESS_LIMIT
	}

	if offset < 0 {
		offset = 0
	}

	newCheckPointsList, err := c.repo.ReadCheckPointsByUserID(ctx, userID, limit, offset)	
	if err != nil {
		return nil, fmt.Errorf("get progress list: %w", err)
	}

	return newCheckPointsList, nil
}

