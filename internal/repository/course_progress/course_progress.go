package courseprogress

import (
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
		return domain.CheckPoint{}, nil //TODO ошибка создания
	}

	return toDomainModel(cp), nil
}

func (r *progressRepo) ReadCheckPointsByUserID(
	ctx context.Context, 
	userID uuid.UUID, 
	limit, 
	offset int,
) ([]domain.CheckPoint, error) {

	var checkPoints []GormCheckPoint
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).
		Limit(limit).
		Offset(offset).
		Find(&checkPoints).
		Error; err != nil {
			return nil, nil //TODO ошибка
	}

	dCheckPoints := make([]domain.CheckPoint, len(checkPoints))
	for i := 0; i < len(checkPoints); i++ {
		dCheckPoints[i] = toDomainModel(checkPoints[i])
	}

	return dCheckPoints, nil
}
