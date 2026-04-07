package contracts

import (
	"context"
	"github.com/google/uuid"
	"github.com/WebCraftersGH/Education-service/internal/domain"
)

//go:generate mockgen -destination=mocks/mock_ProgressSVC.go -package=mocks . ProgressSVC
type ProgressSVC interface {
	ProgressList(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.CheckPoint, error)
	SetProgress(ctx context.Context, userID uuid.UUID, slug string) (domain.CheckPoint, error)
}
