package problemcatalog

import (
	"github.com/WebCraftersGH/Education-service/internal/domain"
)

func ToProblemModel(p domain.Problem) Problem {
	return Problem{
		ID:         p.ID,
		Name:       p.Name,
		Slug:       p.Slug,
		Difficulty: p.Difficulty,
		Tag:        p.Tag,
		Status:     string(p.Status),
		AuthorID:   p.AuthorID,
		VerifiedAt: p.VerifiedAt,
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
	}
}

func ToProblemDomain(p Problem) domain.Problem {
	status := domain.ProblemStatus(p.Status)
	if !status.IsValid() {
		status = domain.ProblemStatusDraft
	}

	return domain.Problem{
		ID:         p.ID,
		Name:       p.Name,
		Slug:       p.Slug,
		Difficulty: p.Difficulty,
		Tag:        p.Tag,
		Status:     status,
		AuthorID:   p.AuthorID,
		VerifiedAt: p.VerifiedAt,
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
	}
}

func ToProblemDomains(problems []Problem) []domain.Problem {
	result := make([]domain.Problem, 0, len(problems))
	for _, p := range problems {
		result = append(result, ToProblemDomain(p))
	}
	return result
}

func ToProblemContentModel(pc domain.ProblemContent) ProblemContent {
	return ProblemContent{
		ID:             pc.ID,
		ProblemID:      pc.ProblemID,
		DescriptionMD:  pc.DescriptionMD,
		InputFormatMD:  pc.InputFormatMD,
		OutputFormatMD: pc.OutputFormatMD,
		ConstraintsMD:  pc.ConstraintsMD,
		NotesMD:        pc.NotesMD,
		CreatedAt:      pc.CreatedAt,
		UpdatedAt:      pc.UpdatedAt,
	}
}

func ToProblemContentDomain(pc ProblemContent) domain.ProblemContent {
	return domain.ProblemContent{
		ID:             pc.ID,
		ProblemID:      pc.ProblemID,
		DescriptionMD:  pc.DescriptionMD,
		InputFormatMD:  pc.InputFormatMD,
		OutputFormatMD: pc.OutputFormatMD,
		ConstraintsMD:  pc.ConstraintsMD,
		NotesMD:        pc.NotesMD,
		CreatedAt:      pc.CreatedAt,
		UpdatedAt:      pc.UpdatedAt,
	}
}
