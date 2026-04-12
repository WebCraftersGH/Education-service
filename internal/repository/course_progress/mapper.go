package courseprogress

import (
	"github.com/WebCraftersGH/Education-service/internal/domain"
)

func toDomainModel(checkPoint GormCheckPoint) domain.CheckPoint {
	return domain.CheckPoint{
		ID: checkPoint.ID,
		UserID: checkPoint.UserID,
		Slug: checkPoint.Slug,
		CreatedAt: checkPoint.CreatedAt,
		UpdatedAt: checkPoint.UpdatedAt,
	}
}

func toGormModel(checkPoint domain.CheckPoint) GormCheckPoint {
	return GormCheckPoint{
		ID: checkPoint.ID,
		UserID: checkPoint.UserID,
		Slug: checkPoint.Slug,
		CreatedAt: checkPoint.CreatedAt,
		UpdatedAt: checkPoint.UpdatedAt,
	}
}
