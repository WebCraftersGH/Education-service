package problemcatalog

import (
	"context"

	"github.com/WebCraftersGH/Education-service/internal/domain"
	"github.com/WebCraftersGH/Education-service/pgk/logging"
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
	logger := logging.GetLogger()

	if model.ID == uuid.Nil {
		model.ID = uuid.New()
	}

	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		logger.WithError(err).WithFields(map[string]any{
			"id":         model.ID,
			"problem_id": model.ProblemID,
		}).Error("create problem content failed")
		return domain.ProblemContent{}, err
	}

	logger.WithFields(map[string]any{
		"id":         model.ID,
		"problem_id": model.ProblemID,
	}).Info("problem content created")

	return ToProblemContentDomain(model), nil
}

func (r *RepositoryProblemContent) ReadByProblemID(ctx context.Context, problemID uuid.UUID) (domain.ProblemContent, error) {
	if problemID == uuid.Nil {
		return domain.ProblemContent{}, domain.ErrProblemContentProblemIDRequired
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
	logger := logging.GetLogger()

	if pc.ProblemID == uuid.Nil {
		return domain.ProblemContent{}, domain.ErrProblemContentProblemIDRequired
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
		logger.WithError(tx.Error).WithField("problem_id", pc.ProblemID).Error("update problem content failed")
		return domain.ProblemContent{}, tx.Error
	}

	if tx.RowsAffected == 0 {
		return domain.ProblemContent{}, gorm.ErrRecordNotFound
	}

	var updated ProblemContent
	if err := r.db.WithContext(ctx).
		Where("problem_id = ?", pc.ProblemID).
		First(&updated).Error; err != nil {
		logger.WithError(err).WithField("problem_id", pc.ProblemID).Error("read updated problem content failed")
		return domain.ProblemContent{}, err
	}

	logger.WithFields(map[string]any{
		"id":         updated.ID,
		"problem_id": updated.ProblemID,
	}).Info("problem content updated")

	return ToProblemContentDomain(updated), nil
}

func (r *RepositoryProblemContent) DeleteByProblemID(ctx context.Context, problemID uuid.UUID) error {
	logger := logging.GetLogger()

	if problemID == uuid.Nil {
		return domain.ErrProblemContentProblemIDRequired
	}

	tx := r.db.WithContext(ctx).
		Where("problem_id = ?", problemID).
		Delete(&ProblemContent{})

	if tx.Error != nil {
		logger.WithError(tx.Error).WithField("problem_id", problemID).Error("delete problem content failed")
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	logger.WithField("problem_id", problemID).Info("problem content deleted")
	return nil
}
