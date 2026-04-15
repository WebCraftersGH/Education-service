package contracts

import (
	"context"
	"github.com/google/uuid"
	"github.com/WebCraftersGH/Education-service/internal/domain"
)

type ProgressSVC interface {
	ProgressList(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.CheckPoint, error)
	SetProgress(ctx context.Context, userID uuid.UUID, slug string) (domain.CheckPoint, error)
}

//go:generate sh -c "go build -o ./mocks/contracts.a . && mockgen -archive=./mocks/contracts.a -destination=./mocks/mock_ProgressRepo.go -package=mocks github.com/WebCraftersGH/Education-service/internal/contracts ProgressRepo"
type ProgressRepo interface {
	CreateCheckPoint(ctx context.Context, checkPoint domain.CheckPoint) (domain.CheckPoint, error)
	ReadCheckPointsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.CheckPoint, error)
}

