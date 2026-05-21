package problemcatalog

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/WebCraftersGH/Education-service/internal/domain"
	"github.com/google/uuid"
)

type fakeProblemContentRepo struct {
	content domain.ProblemContent
}

func (r fakeProblemContentRepo) Create(context.Context, domain.ProblemContent) (domain.ProblemContent, error) {
	return domain.ProblemContent{}, nil
}

func (r fakeProblemContentRepo) ReadByProblemID(context.Context, uuid.UUID) (domain.ProblemContent, error) {
	return r.content, nil
}

func (r fakeProblemContentRepo) Update(context.Context, domain.ProblemContent) (domain.ProblemContent, error) {
	return domain.ProblemContent{}, nil
}

func (r fakeProblemContentRepo) DeleteByProblemID(context.Context, uuid.UUID) error {
	return nil
}

func TestSolveComparesActualGraphWithExpectedGraph(t *testing.T) {
	problemID := uuid.New()
	uc := NewProblemContentUseCase(fakeProblemContentRepo{
		content: domain.ProblemContent{
			ProblemID:     problemID,
			ExpectedGraph: json.RawMessage(`{"nodes":[{"id":"a"}],"edges":[]}`),
		},
	})

	solved, err := uc.Solve(context.Background(), problemID, json.RawMessage(`{"edges":[],"nodes":[{"id":"a"}]}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !solved {
		t.Fatal("expected graph to be solved")
	}
}

func TestSolveReturnsFalseForDifferentGraph(t *testing.T) {
	problemID := uuid.New()
	uc := NewProblemContentUseCase(fakeProblemContentRepo{
		content: domain.ProblemContent{
			ProblemID:     problemID,
			ExpectedGraph: json.RawMessage(`{"nodes":[{"id":"a"}],"edges":[]}`),
		},
	})

	solved, err := uc.Solve(context.Background(), problemID, json.RawMessage(`{"nodes":[{"id":"b"}],"edges":[]}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if solved {
		t.Fatal("expected graph to be unsolved")
	}
}

func TestSolveRejectsInvalidActualGraph(t *testing.T) {
	uc := NewProblemContentUseCase(fakeProblemContentRepo{})

	_, err := uc.Solve(context.Background(), uuid.New(), json.RawMessage(`{`))
	if !errors.Is(err, domain.ErrProblemContentActualGraphInvalid) {
		t.Fatalf("expected actual graph validation error, got %v", err)
	}
}
