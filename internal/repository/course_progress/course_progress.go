package courseprogress

import (
	"context"
	"errors"

	"github.com/WebCraftersGH/Education-service/internal/domain"
	"github.com/WebCraftersGH/Education-service/internal/logctx"
	"github.com/WebCraftersGH/Education-service/pgk/logging"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type progressRepo struct {
	db *gorm.DB
	logger logging.Logger	
}

func NewProgressRepo(db *gorm.DB, logger logging.Logger) *progressRepo {
	return &progressRepo{db:db, logger: logger}
}

func (r *progressRepo) CreateCheckPoint(
	ctx context.Context, 
	checkPoint domain.CheckPoint,
) (domain.CheckPoint, error) {

	logger := logctx.WithContext(ctx, r.logger).WithFields(map[string]any{
		"user_id": checkPoint.UserID.String(),
		"checkpoint_slug": checkPoint.Slug,
		"repo_method": "CreateCheckPoint",
	})

	cp := toGormModel(checkPoint)
	if err := r.db.WithContext(ctx).Create(&cp).Error; err != nil {

		switch {
		case errors.Is(err, gorm.ErrDuplicatedKey):
			logger.WithError(err).Info("the record already exists")
			return domain.CheckPoint{}, ErrDuplicateRecord
		default:
			logger.WithError(err).Error("create checkpoint failed")
			return domain.CheckPoint{}, ErrInternal
		}
	}

	return toDomainModel(cp), nil
}

func (r *progressRepo) ReadCheckPointsByUserID(
	ctx context.Context, 
	userID uuid.UUID, 
	limit, 
	offset int,
) ([]domain.CheckPoint, error) {
	
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	if limit > 100 { //TODO Это на первое время, потом надо менять.
		limit = 100
	}
	
	logger := logctx.WithContext(ctx, r.logger).WithFields(map[string]any{
		"user_id": userID.String(),
		"limit": limit,
		"offset": offset,
		"repo_method": "ReadCheckPointsByUserID",
	})

	var checkPoints []GormCheckPoint
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).
		Limit(limit).
		Offset(offset).
		Find(&checkPoints).
		Error; err != nil {
			switch { //TODO Лишняя обертка, но пока оставлю
				default:
					logger.WithError(err).Error("read checkpoints by user id failed")
					return nil, ErrInternal
			}
	}

	dCheckPoints := make([]domain.CheckPoint, len(checkPoints))
	for i := 0; i < len(checkPoints); i++ {
		dCheckPoints[i] = toDomainModel(checkPoints[i])
	}

	return dCheckPoints, nil
}
