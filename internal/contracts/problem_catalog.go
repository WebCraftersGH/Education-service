package contracts

import (
	"context"
	"encoding/json"
	"github.com/WebCraftersGH/Education-service/internal/domain"
	"github.com/google/uuid"
)

type ProblemSVC interface {
	Create(ctx context.Context, p domain.Problem) (domain.Problem, error)
	ReadBySlug(ctx context.Context, pSlug string) (domain.Problem, error)
	ReadByID(ctx context.Context, problemID uuid.UUID) (domain.Problem, error)
	Update(ctx context.Context, p domain.Problem) (domain.Problem, error)
	DeleteBySlug(ctx context.Context, pSlug string) error
	List(ctx context.Context, filter domain.ProblemFilter) ([]domain.Problem, error)
}

type ProblemContentSVC interface {
	Create(ctx context.Context, pc domain.ProblemContent) (domain.ProblemContent, error)
	ReadByProblemID(ctx context.Context, problemID uuid.UUID) (domain.ProblemContent, error)
	Update(ctx context.Context, pc domain.ProblemContent) (domain.ProblemContent, error)
	DeleteByProblemID(ctx context.Context, problemID uuid.UUID) error
	Solve(ctx context.Context, problemID uuid.UUID, actualGraph json.RawMessage) (bool, error)
}

type ProblemRepository interface {
	Create(ctx context.Context, p domain.Problem) (domain.Problem, error)
	ReadBySlug(ctx context.Context, pSlug string) (domain.Problem, error)
	ReadByID(ctx context.Context, problemID uuid.UUID) (domain.Problem, error)
	Update(ctx context.Context, p domain.Problem) (domain.Problem, error)
	DeleteBySlug(ctx context.Context, pSlug string) error
	List(ctx context.Context, filter domain.ProblemFilter) ([]domain.Problem, error)
}

type ProblemContentRepository interface {
	Create(ctx context.Context, pc domain.ProblemContent) (domain.ProblemContent, error)
	ReadByProblemID(ctx context.Context, problemID uuid.UUID) (domain.ProblemContent, error)
	Update(ctx context.Context, pc domain.ProblemContent) (domain.ProblemContent, error)
	DeleteByProblemID(ctx context.Context, problemID uuid.UUID) error
}
