package problemcatalog

import (
	"context"
	"strings"

	"github.com/WebCraftersGH/Education-service/internal/contracts"
	"github.com/WebCraftersGH/Education-service/internal/domain"
	"github.com/google/uuid"
)

type ProblemUseCase struct {
	repo contracts.ProblemRepository
}

func NewProblemUseCase(repo contracts.ProblemRepository) *ProblemUseCase {
	return &ProblemUseCase{repo: repo}
}

func (uc *ProblemUseCase) Create(ctx context.Context, p domain.Problem) (domain.Problem, error) {
	p.Name = strings.TrimSpace(p.Name)
	p.Slug = strings.TrimSpace(strings.ToLower(p.Slug))
	p.Difficulty = strings.TrimSpace(strings.ToLower(p.Difficulty))
	p.Tag = strings.TrimSpace(strings.ToLower(p.Tag))

	if p.Name == "" {
		return domain.Problem{}, domain.ErrProblemNameRequired
	}

	if p.Slug == "" {
		return domain.Problem{}, domain.ErrProblemSlugRequired
	}

	if p.Difficulty == "" {
		return domain.Problem{}, domain.ErrProblemDifficultyRequired
	}

	if p.AuthorID == uuid.Nil {
		return domain.Problem{}, domain.ErrProblemAuthorIDRequired
	}

	if !p.Status.IsValid() {
		p.Status = domain.ProblemStatusDraft
	}

	return uc.repo.Create(ctx, p)
}

func (uc *ProblemUseCase) ReadBySlug(ctx context.Context, pSlug string) (domain.Problem, error) {
	pSlug = strings.TrimSpace(strings.ToLower(pSlug))
	if pSlug == "" {
		return domain.Problem{}, domain.ErrProblemSlugRequired
	}

	return uc.repo.ReadBySlug(ctx, pSlug)
}

func (uc *ProblemUseCase) Update(ctx context.Context, p domain.Problem) (domain.Problem, error) {
	if p.ID == uuid.Nil {
		return domain.Problem{}, domain.ErrProblemIDRequired
	}

	p.Name = strings.TrimSpace(p.Name)
	p.Slug = strings.TrimSpace(strings.ToLower(p.Slug))
	p.Difficulty = strings.TrimSpace(strings.ToLower(p.Difficulty))
	p.Tag = strings.TrimSpace(strings.ToLower(p.Tag))

	if p.Name == "" {
		return domain.Problem{}, domain.ErrProblemNameRequired
	}

	if p.Slug == "" {
		return domain.Problem{}, domain.ErrProblemSlugRequired
	}

	if p.Difficulty == "" {
		return domain.Problem{}, domain.ErrProblemDifficultyRequired
	}

	if p.AuthorID == uuid.Nil {
		return domain.Problem{}, domain.ErrProblemAuthorIDRequired
	}

	if !p.Status.IsValid() {
		p.Status = domain.ProblemStatusDraft
	}

	return uc.repo.Update(ctx, p)
}

func (uc *ProblemUseCase) DeleteBySlug(ctx context.Context, pSlug string) error {
	pSlug = strings.TrimSpace(strings.ToLower(pSlug))
	if pSlug == "" {
		return domain.ErrProblemSlugRequired
	}

	return uc.repo.DeleteBySlug(ctx, pSlug)
}

func (uc *ProblemUseCase) List(ctx context.Context, filter domain.ProblemFilter) ([]domain.Problem, error) {
	filter.Tag = strings.TrimSpace(strings.ToLower(filter.Tag))
	filter.Difficulty = strings.TrimSpace(strings.ToLower(filter.Difficulty))

	if filter.Limit <= 0 {
		filter.Limit = 20
	}

	if filter.Offset < 0 {
		filter.Offset = 0
	}

	return uc.repo.List(ctx, filter)
}
