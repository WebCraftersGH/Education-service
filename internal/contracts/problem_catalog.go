package contracts

import (
	"context"
	"github.com/WebCraftersGH/Education-service/internal/domain"
	"github.com/google/uuid"
)

type ProblemSVC interface {
	Create(ctx context.Context, p domain.Problem) (domain.Problem, error)
	ReadBySlug(ctx context.Context, pSlug string) (domain.Problem, error)
	Update(ctx context.Context, p domain.Problem) (domain.Problem, error)
	DeleteBySlug(ctx context.Context, pSlug string) error
	List(ctx context.Context, filter domain.ProblemFilter) ([]domain.Problem, error)
}

type ProblemContentSVC interface {
	Create(ctx context.Context, pc domain.ProblemContent) (domain.ProblemContent, error)
	ReadByProblemID(ctx context.Context, problemID uuid.UUID) (domain.ProblemContent, error)
	Update(ctx context.Context, pc domain.ProblemContent) (domain.ProblemContent, error)
	DeleteByProblemID(ctx context.Context, problemID uuid.UUID) error
}

type ProblemRepository interface {
	Create(ctx context.Context, p domain.Problem) (domain.Problem, error)
	ReadBySlug(ctx context.Context, pSlug string) (domain.Problem, error)
	Update(ctx context.Context, p domain.Problem) (domain.Problem, error)
	DeleteBySlug(ctx context.Context, pSlug string) error
	List(ctx context.Context, filter domain.ProblemFilter) ([]domain.Problem, error)
}
