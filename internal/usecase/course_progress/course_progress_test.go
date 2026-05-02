package courseprogress

import (
	"context"
	"errors"
	"testing"

	"github.com/WebCraftersGH/Education-service/internal/apperrors"
	"github.com/WebCraftersGH/Education-service/internal/contracts/mocks"
	"github.com/WebCraftersGH/Education-service/internal/domain"
	"github.com/WebCraftersGH/Education-service/pkg/logging"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestSetProgress_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger, closer, _ := logging.New("INFO")
	defer closer.Close()

	ctx := context.Background()
	repo := mocks.NewMockProgressRepo(ctrl)
	svc := NewCourseProgress(repo, logger)

	userID := uuid.New()
	inputSlug := "    New-slug     "
	expectedSlug := "New-slug"

	repo.EXPECT().CreateCheckPoint(ctx, domain.CheckPoint{
		UserID: userID,
		Slug:   expectedSlug,
	}).Return(domain.CheckPoint{
		UserID: userID,
		Slug:   expectedSlug,
	}, nil).Times(1)

	checkpoint, err := svc.SetProgress(ctx, userID, inputSlug)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if checkpoint.Slug != expectedSlug {
		t.Fatalf("got slug %s, want %s", checkpoint.Slug, expectedSlug)
	}
}

func TestSetProgress_DuplicateRecord(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger, closer, _ := logging.New("INFO")
	defer closer.Close()

	ctx := context.Background()
	repo := mocks.NewMockProgressRepo(ctrl)
	svc := NewCourseProgress(repo, logger)

	userID := uuid.New()
	slug := "New-slug"

	repo.EXPECT().CreateCheckPoint(ctx, domain.CheckPoint{
		UserID: userID,
		Slug:   slug,
	}).Return(domain.CheckPoint{}, apperrors.ErrDuplicateRecord).
		Times(1)

	_, err := svc.SetProgress(ctx, userID, slug)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if !errors.Is(err, apperrors.ErrDuplicateRecord) {
		t.Fatalf("expected ErrDuplicateRecord, got %v", err)
	}
}

func TestSetProgress_EmptySlug(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger, closer, _ := logging.New("INFO")
	defer closer.Close()

	ctx := context.Background()
	repo := mocks.NewMockProgressRepo(ctrl)
	svc := NewCourseProgress(repo, logger)

	userID := uuid.New()
	slug := "  "

	_, err := svc.SetProgress(ctx, userID, slug)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if !errors.Is(err, apperrors.ErrInvalidArgument) {
		t.Fatalf("expected ErrInvalidArgument, got %v", err)
	}
}

func TestSetProgress_EmptyUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger, closer, _ := logging.New("INFO")
	defer closer.Close()

	ctx := context.Background()
	repo := mocks.NewMockProgressRepo(ctrl)
	svc := NewCourseProgress(repo, logger)

	userID := uuid.Nil
	slug := "New-slug"

	_, err := svc.SetProgress(ctx, userID, slug)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if !errors.Is(err, apperrors.ErrInvalidArgument) {
		t.Fatalf("expected ErrInvalidArgument, got %v", err)
	}

}

func TestSetProgress_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger, closer, _ := logging.New("INFO")
	defer closer.Close()

	ctx := context.Background()
	repo := mocks.NewMockProgressRepo(ctrl)
	svc := NewCourseProgress(repo, logger)

	userID := uuid.New()
	slug := "New-slug"

	repo.EXPECT().
		CreateCheckPoint(ctx, domain.CheckPoint{
			UserID: userID,
			Slug:   slug,
		}).
		Return(domain.CheckPoint{}, apperrors.ErrInternal).
		Times(1)

	_, err := svc.SetProgress(ctx, userID, slug)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, apperrors.ErrInternal) {
		t.Fatalf("expected ErrInternal, got %v", err)
	}
}

func TestProgressList_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger, closer, _ := logging.New("INFO")
	defer closer.Close()

	ctx := context.Background()
	repo := mocks.NewMockProgressRepo(ctrl)
	svc := NewCourseProgress(repo, logger)

	userID := uuid.New()
	expected := []domain.CheckPoint{
		{
			ID:     uuid.New(),
			UserID: userID,
			Slug:   "module-3",
		},
		{
			ID:     uuid.New(),
			UserID: userID,
			Slug:   "module-2",
		},
	}

	repo.EXPECT().
		ReadCheckPointsByUserID(ctx, userID, 5, 0).
		Return(expected, nil).
		Times(1)

	got, err := svc.ProgressList(ctx, userID, 5, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != len(expected) {
		t.Fatalf("got len %d, want %d", len(got), len(expected))
	}

	for i := range got {
		if got[i].ID != expected[i].ID {
			t.Fatalf("got id %v, want %v", got[i].ID, expected[i].ID)
		}
		if got[i].UserID != expected[i].UserID {
			t.Fatalf("got userID %v, want %v", got[i].UserID, expected[i].UserID)
		}
		if got[i].Slug != expected[i].Slug {
			t.Fatalf("got slug %q, want %q", got[i].Slug, expected[i].Slug)
		}
	}
}

func TestProgressList_DefaultLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger, closer, _ := logging.New("INFO")
	defer closer.Close()

	ctx := context.Background()
	repo := mocks.NewMockProgressRepo(ctrl)
	svc := NewCourseProgress(repo, logger)

	userID := uuid.New()

	repo.EXPECT().
		ReadCheckPointsByUserID(ctx, userID, BASE_GET_PROGRESS_LIMIT, 0).
		Return([]domain.CheckPoint{}, nil).
		Times(1)

	_, err := svc.ProgressList(ctx, userID, 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestProgressList_MaxLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger, closer, _ := logging.New("INFO")
	defer closer.Close()

	ctx := context.Background()
	repo := mocks.NewMockProgressRepo(ctrl)
	svc := NewCourseProgress(repo, logger)

	userID := uuid.New()

	repo.EXPECT().
		ReadCheckPointsByUserID(ctx, userID, MAX_PROGRESS_LIMIT, 0).
		Return([]domain.CheckPoint{}, nil).
		Times(1)

	_, err := svc.ProgressList(ctx, userID, MAX_PROGRESS_LIMIT+1, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestProgressList_NegativeOffset(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger, closer, _ := logging.New("INFO")
	defer closer.Close()

	ctx := context.Background()
	repo := mocks.NewMockProgressRepo(ctrl)
	svc := NewCourseProgress(repo, logger)

	userID := uuid.New()

	repo.EXPECT().
		ReadCheckPointsByUserID(ctx, userID, 10, 0).
		Return([]domain.CheckPoint{}, nil).
		Times(1)

	_, err := svc.ProgressList(ctx, userID, 10, -5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestProgressList_EmptyUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger, closer, _ := logging.New("INFO")
	defer closer.Close()

	ctx := context.Background()
	repo := mocks.NewMockProgressRepo(ctrl)
	svc := NewCourseProgress(repo, logger)

	_, err := svc.ProgressList(ctx, uuid.Nil, 10, 0)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, apperrors.ErrInvalidArgument) {
		t.Fatalf("expected ErrInvalidArgument, got %v", err)
	}
}

func TestProgressList_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger, closer, _ := logging.New("INFO")
	defer closer.Close()

	ctx := context.Background()
	repo := mocks.NewMockProgressRepo(ctrl)
	svc := NewCourseProgress(repo, logger)

	userID := uuid.New()

	repo.EXPECT().
		ReadCheckPointsByUserID(ctx, userID, 10, 0).
		Return(nil, apperrors.ErrInternal).
		Times(1)

	_, err := svc.ProgressList(ctx, userID, 10, 0)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, apperrors.ErrInternal) {
		t.Fatalf("expected ErrInternal, got %v", err)
	}
}

func TestProgressList_NegativeLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger, closer, _ := logging.New("INFO")
	defer closer.Close()

	ctx := context.Background()
	repo := mocks.NewMockProgressRepo(ctrl)
	svc := NewCourseProgress(repo, logger)

	userID := uuid.New()

	repo.EXPECT().
		ReadCheckPointsByUserID(ctx, userID, BASE_GET_PROGRESS_LIMIT, 0).
		Return([]domain.CheckPoint{}, nil).
		Times(1)

	_, err := svc.ProgressList(ctx, userID, -5, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestProgressList_PositiveOffset(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger, closer, _ := logging.New("INFO")
	defer closer.Close()

	ctx := context.Background()
	repo := mocks.NewMockProgressRepo(ctrl)
	svc := NewCourseProgress(repo, logger)

	userID := uuid.New()

	repo.EXPECT().
		ReadCheckPointsByUserID(ctx, userID, 10, 3).
		Return([]domain.CheckPoint{}, nil).
		Times(1)

	_, err := svc.ProgressList(ctx, userID, 10, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
