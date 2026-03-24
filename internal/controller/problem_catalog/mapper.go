package problemcatalog

import (
	"github.com/WebCraftersGH/Education-service/internal/domain"
)

func toProblemResponse(p domain.Problem) problemResponse {
	return problemResponse{
		ID:         p.ID,
		Name:       p.Name,
		Slug:       p.Slug,
		Difficulty: p.Difficulty,
		Tag:        p.Tag,
		Status:     p.Status,
		AuthorID:   p.AuthorID,
		VerifiedAt: p.VerifiedAt,
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
	}
}

func toProblemsResponse(items []domain.Problem) []problemResponse {
	result := make([]problemResponse, 0, len(items))
	for _, item := range items {
		result = append(result, toProblemResponse(item))
	}
	return result
}

func toProblemContentResponse(pc domain.ProblemContent) problemContentResponse {
	return problemContentResponse{
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
