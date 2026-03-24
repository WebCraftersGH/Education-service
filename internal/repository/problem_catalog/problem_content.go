package problemcatalog

import (
	"context"

	"github.com/WebCraftersGH/Education-service/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RepositoryProblemContent struct {
	db *gorm.DB
}

func NewRepositoryProblemContent(db *gorm.DB) *RepositoryProblemContent {
	return &RepositoryProblemContent{db: db}
}

func (r *RepositoryProblemContent) Create(ctx context.Context, pc domain.ProblemContent) (domain.ProblemContent, error) {
	model := ToProblemContentModel(pc)

	if model.ID == uuid.Nil {
		model.ID = uuid.New()
	}

	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return domain.ProblemContent{}, err
	}

	return ToProblemContentDomain(model), nil
}

func (r *RepositoryProblemContent) ReadByProblemID(ctx context.Context, problemID uuid.UUID) (domain.ProblemContent, error) {
	if problemID == uuid.Nil {
		return domain.ProblemContent{}, domain.ErrProblemIDRequired
	}

	var model ProblemContent

	if err := r.db.WithContext(ctx).
		Where("problem_id = ?", problemID).
		First(&model).Error; err != nil {
		return domain.ProblemContent{}, err
	}

	return ToProblemContentDomain(model), nil
}

func (r *RepositoryProblemContent) Update(ctx context.Context, pc domain.ProblemContent) (domain.ProblemContent, error) {
	if pc.ProblemID == uuid.Nil {
		return domain.ProblemContent{}, domain.ErrProblemIDRequired
	}

	model := ToProblemContentModel(pc)

	tx := r.db.WithContext(ctx).
		Model(&ProblemContent{}).
		Where("problem_id = ?", pc.ProblemID).
		Updates(map[string]any{
			"description_md":   model.DescriptionMD,
			"input_format_md":  model.InputFormatMD,
			"output_format_md": model.OutputFormatMD,
			"constraints_md":   model.ConstraintsMD,
			"notes_md":         model.NotesMD,
		})

	if tx.Error != nil {
		return domain.ProblemContent{}, tx.Error
	}

	if tx.RowsAffected == 0 {
		return domain.ProblemContent{}, gorm.ErrRecordNotFound
	}

	var updated ProblemContent
	if err := r.db.WithContext(ctx).
		Where("problem_id = ?", pc.ProblemID).
		First(&updated).Error; err != nil {
		return domain.ProblemContent{}, err
	}

	return ToProblemContentDomain(updated), nil
}

func (r *RepositoryProblemContent) DeleteByProblemID(ctx context.Context, problemID uuid.UUID) error {
	if problemID == uuid.Nil {
		return domain.ErrProblemIDRequired
	}

	tx := r.db.WithContext(ctx).
		Where("problem_id = ?", problemID).
		Delete(&ProblemContent{})

	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
