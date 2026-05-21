package problemcatalog

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"

	"github.com/WebCraftersGH/Education-service/internal/contracts"
	"github.com/WebCraftersGH/Education-service/internal/domain"
	"github.com/google/uuid"
)

type ProblemContentUseCase struct {
	repo contracts.ProblemContentRepository
}

func NewProblemContentUseCase(repo contracts.ProblemContentRepository) *ProblemContentUseCase {
	return &ProblemContentUseCase{repo: repo}
}

func (uc *ProblemContentUseCase) Create(ctx context.Context, pc domain.ProblemContent) (domain.ProblemContent, error) {
	if pc.ProblemID == uuid.Nil {
		return domain.ProblemContent{}, domain.ErrProblemContentProblemIDRequired
	}

	if pc.AuthorID == uuid.Nil {
		return domain.ProblemContent{}, domain.ErrProblemAuthorIDRequired
	}

	pc.FullText = strings.TrimSpace(pc.FullText)

	if pc.FullText == "" {
		return domain.ProblemContent{}, domain.ErrProblemContentFullTextRequired
	}

	if !json.Valid(pc.ActualGraph) {
		return domain.ProblemContent{}, domain.ErrProblemContentActualGraphInvalid
	}

	if !json.Valid(pc.ExpectedGraph) {
		return domain.ProblemContent{}, domain.ErrProblemContentExpectedGraphInvalid
	}

	return uc.repo.Create(ctx, pc)
}

func (uc *ProblemContentUseCase) ReadByProblemID(ctx context.Context, problemID uuid.UUID) (domain.ProblemContent, error) {
	if problemID == uuid.Nil {
		return domain.ProblemContent{}, domain.ErrProblemContentProblemIDRequired
	}

	return uc.repo.ReadByProblemID(ctx, problemID)
}

func (uc *ProblemContentUseCase) Update(ctx context.Context, pc domain.ProblemContent) (domain.ProblemContent, error) {
	if pc.ID == uuid.Nil {
		return domain.ProblemContent{}, domain.ErrProblemContentIDRequired
	}

	if pc.ProblemID == uuid.Nil {
		return domain.ProblemContent{}, domain.ErrProblemContentProblemIDRequired
	}

	if pc.AuthorID == uuid.Nil {
		return domain.ProblemContent{}, domain.ErrProblemAuthorIDRequired
	}

	pc.FullText = strings.TrimSpace(pc.FullText)

	if pc.FullText == "" {
		return domain.ProblemContent{}, domain.ErrProblemContentFullTextRequired
	}

	if !json.Valid(pc.ActualGraph) {
		return domain.ProblemContent{}, domain.ErrProblemContentActualGraphInvalid
	}

	if !json.Valid(pc.ExpectedGraph) {
		return domain.ProblemContent{}, domain.ErrProblemContentExpectedGraphInvalid
	}

	return uc.repo.Update(ctx, pc)
}

func (uc *ProblemContentUseCase) DeleteByProblemID(ctx context.Context, problemID uuid.UUID) error {
	if problemID == uuid.Nil {
		return domain.ErrProblemContentProblemIDRequired
	}

	return uc.repo.DeleteByProblemID(ctx, problemID)
}

func (uc *ProblemContentUseCase) Solve(ctx context.Context, problemID uuid.UUID, actualGraph json.RawMessage) (bool, error) {
	if problemID == uuid.Nil {
		return false, domain.ErrProblemContentProblemIDRequired
	}

	if !json.Valid(actualGraph) {
		return false, domain.ErrProblemContentActualGraphInvalid
	}

	pc, err := uc.repo.ReadByProblemID(ctx, problemID)
	if err != nil {
		return false, err
	}

	var actual any
	if err := json.Unmarshal(actualGraph, &actual); err != nil {
		return false, domain.ErrProblemContentActualGraphInvalid
	}

	var expected any
	if err := json.Unmarshal(pc.ExpectedGraph, &expected); err != nil {
		return false, domain.ErrProblemContentExpectedGraphInvalid
	}

	return reflect.DeepEqual(actual, expected), nil
}
