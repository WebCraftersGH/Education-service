package courseprogress

import (
	"fmt"
	"errors"
	"context"

	"github.com/WebCraftersGH/Education-service/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type progressRepo struct {
	db *gorm.DB
}

func NewProgressRepo(db *gorm.DB) *progressRepo {
	return &progressRepo{db:db}
}

func (r *progressRepo) CreateCheckPoint(
	ctx context.Context, 
	checkPoint domain.CheckPoint,
) (domain.CheckPoint, error) {

	cp := toGormModel(checkPoint)
	if err := r.db.WithContext(ctx).Create(&cp).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return domain.CheckPoint{}, fmt.Errorf("gormrepo: create checkpoint: %w", ErrDuplicateRecord)
		}
		return domain.CheckPoint{}, fmt.Errorf("gormrepo: %v: %w", err, ErrInternal)
	}

	return toDomainModel(cp), nil
}

func (r *progressRepo) ReadCheckPointsByUserID(
	ctx context.Context, 
	userID uuid.UUID, 
	limit, 
	offset int,
) ([]domain.CheckPoint, error) {

	if limit < 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	if limit > 100 { //TODO Это на первое время, потом надо менять.
		limit = 100
	}

	var checkPoints []GormCheckPoint
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).
		Limit(limit).
		Offset(offset).
		Find(&checkPoints).
		Error; err != nil {
			return nil, fmt.Errorf("gormrepo: %v: %w", err, ErrInternal)
	}

	dCheckPoints := make([]domain.CheckPoint, len(checkPoints))
	for i := 0; i < len(checkPoints); i++ {
		dCheckPoints[i] = toDomainModel(checkPoints[i])
	}

	return dCheckPoints, nil
}
