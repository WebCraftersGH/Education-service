package problemcatalog

import (
	"context"
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

	pc.DescriptionMD = strings.TrimSpace(pc.DescriptionMD)
	pc.InputFormatMD = strings.TrimSpace(pc.InputFormatMD)
	pc.OutputFormatMD = strings.TrimSpace(pc.OutputFormatMD)
	pc.ConstraintsMD = strings.TrimSpace(pc.ConstraintsMD)
	pc.NotesMD = strings.TrimSpace(pc.NotesMD)

	if pc.DescriptionMD == "" {
		return domain.ProblemContent{}, domain.ErrProblemDescriptionRequired
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

	pc.DescriptionMD = strings.TrimSpace(pc.DescriptionMD)
	pc.InputFormatMD = strings.TrimSpace(pc.InputFormatMD)
	pc.OutputFormatMD = strings.TrimSpace(pc.OutputFormatMD)
	pc.ConstraintsMD = strings.TrimSpace(pc.ConstraintsMD)
	pc.NotesMD = strings.TrimSpace(pc.NotesMD)

	if pc.DescriptionMD == "" {
		return domain.ProblemContent{}, domain.ErrProblemDescriptionRequired
	}

	return uc.repo.Update(ctx, pc)
}

func (uc *ProblemContentUseCase) DeleteByProblemID(ctx context.Context, problemID uuid.UUID) error {
	if problemID == uuid.Nil {
		return domain.ErrProblemContentProblemIDRequired
	}

	return uc.repo.DeleteByProblemID(ctx, problemID)
}
