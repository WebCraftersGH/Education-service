package handlers

import (
	"github.com/WebCraftersGH/Education-service/internal/domain"
)

func toProgressResponse(checkpoint domain.CheckPoint) ProgressResponse {
	return ProgressResponse{
		ID:        checkpoint.ID,
		UserID:    checkpoint.UserID,
		Slug:      checkpoint.Slug,
		CreatedAt: checkpoint.CreatedAt,
		UpdatedAt: checkpoint.UpdatedAt,
	}
}

func toProgressResponseList(checkpoints []domain.CheckPoint) ProgressResponseList {
	list := make([]ProgressResponse, len(checkpoints))
	for i, checkpoint := range checkpoints {
		ck := toProgressResponse(checkpoint)
		list[i] = ck
	}
	return ProgressResponseList{ProgressList: list}
}

func toProblemResponse(problem domain.Problem) ProblemResponse {
	return ProblemResponse{
		ID:         problem.ID,
		Name:       problem.Name,
		Slug:       problem.Slug,
		Difficulty: problem.Difficulty,
		Tag:        problem.Tag,
		Status:     string(problem.Status),
		AuthorID:   problem.AuthorID,
		VerifiedAt: problem.VerifiedAt,
		CreatedAt:  problem.CreatedAt,
		UpdatedAt:  problem.UpdatedAt,
	}
}

func toProblemResponseList(problems []domain.Problem) ProblemResponseList {
	list := make([]ProblemResponse, len(problems))
	for i, problem := range problems {
		list[i] = toProblemResponse(problem)
	}

	return ProblemResponseList{Problems: list}
}

func toProblemContentResponse(problemContent domain.ProblemContent) ProblemContentResponse {
	return ProblemContentResponse{
		ID:            problemContent.ID,
		ProblemID:     problemContent.ProblemID,
		AuthorID:      problemContent.AuthorID,
		ActualGraph:   problemContent.ActualGraph,
		ExpectedGraph: problemContent.ExpectedGraph,
		FullText:      problemContent.FullText,
		CreatedAt:     problemContent.CreatedAt,
		UpdatedAt:     problemContent.UpdatedAt,
	}
}
