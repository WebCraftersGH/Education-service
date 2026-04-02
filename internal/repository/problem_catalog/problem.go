package problemcatalog

import (
	"context"
	"errors"
	"strings"

	"github.com/WebCraftersGH/Education-service/internal/contracts"
	"github.com/WebCraftersGH/Education-service/internal/domain"
	"github.com/WebCraftersGH/Education-service/pgk/logging"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) contracts.ProblemRepository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, p domain.Problem) (domain.Problem, error) {
	model := ToProblemModel(p)
	logger := logging.GetLogger()

	if model.ID == uuid.Nil {
		model.ID = uuid.New()
	}

	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		logger.WithError(err).WithFields(map[string]any{
			"slug": model.Slug,
			"id":   model.ID,
		}).Error("create problem failed")
		return domain.Problem{}, err
	}

	logger.WithFields(map[string]any{
		"slug": model.Slug,
		"id":   model.ID,
	}).Info("problem created")

	return ToProblemDomain(model), nil
}

func (r *Repository) ReadBySlug(ctx context.Context, pSlug string) (domain.Problem, error) {
	var model Problem

	if err := r.db.WithContext(ctx).
		Where("slug = ?", strings.TrimSpace(strings.ToLower(pSlug))).
		First(&model).Error; err != nil {
		return domain.Problem{}, err
	}

	return ToProblemDomain(model), nil
}

func (r *Repository) Read(ctx context.Context, problemID uuid.UUID) (domain.Problem, error) {
	var model Problem

	if err := r.db.WithContext(ctx).
		Where("id = ?", problemID).
		First(&model).Error; err != nil {
		return domain.Problem{}, err
	}

	return ToProblemDomain(model), nil
}

func (r *Repository) Update(ctx context.Context, p domain.Problem) (domain.Problem, error) {
	logger := logging.GetLogger()

	if p.ID == uuid.Nil {
		return domain.Problem{}, errors.New("problem id is required")
	}

	model := ToProblemModel(p)

	tx := r.db.WithContext(ctx).
		Model(&Problem{}).
		Where("id = ?", p.ID).
		Updates(map[string]any{
			"name":        model.Name,
			"slug":        model.Slug,
			"difficulty":  model.Difficulty,
			"tag":         model.Tag,
			"status":      model.Status,
			"author_id":   model.AuthorID,
			"verified_at": model.VerifiedAt,
		})

	if tx.Error != nil {
		logger.WithError(tx.Error).WithField("id", p.ID).Error("update problem failed")
		return domain.Problem{}, tx.Error
	}

	if tx.RowsAffected == 0 {
		return domain.Problem{}, gorm.ErrRecordNotFound
	}

	var updated Problem
	if err := r.db.WithContext(ctx).
		Where("id = ?", p.ID).
		First(&updated).Error; err != nil {
		logger.WithError(err).WithField("id", p.ID).Error("read updated problem failed")
		return domain.Problem{}, err
	}

	logger.WithFields(map[string]any{
		"id":   p.ID,
		"slug": updated.Slug,
	}).Info("problem updated")

	return ToProblemDomain(updated), nil
}

func (r *Repository) DeleteBySlug(ctx context.Context, pSlug string) error {
	logger := logging.GetLogger()
	tx := r.db.WithContext(ctx).
		Where("slug = ?", strings.TrimSpace(strings.ToLower(pSlug))).
		Delete(&Problem{})

	if tx.Error != nil {
		logger.WithError(tx.Error).WithField("slug", pSlug).Error("delete problem by slug failed")
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	logger.WithField("slug", pSlug).Info("problem deleted")
	return nil
}

func (r *Repository) Delete(ctx context.Context, problemID uuid.UUID) error {
	tx := r.db.WithContext(ctx).
		Where("id = ?", problemID).
		Delete(&Problem{})

	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *Repository) List(ctx context.Context, filter domain.ProblemFilter) ([]domain.Problem, error) {
	var models []Problem

	query := r.db.WithContext(ctx).Model(&Problem{})

	if filter.Tag != "" {
		query = query.Where("tag = ?", filter.Tag)
	}

	if filter.Difficulty != "" {
		query = query.Where("difficulty = ?", filter.Difficulty)
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}

	offset := filter.Offset
	if offset < 0 {
		offset = 0
	}

	if err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&models).Error; err != nil {
		return nil, err
	}

	return ToProblemDomains(models), nil
}
