package courseprogress

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormCheckPoint struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:checkpoints_user_id_slug_key"`
	Slug string `gorm:"varchar(255);not null;uniqueIndex:checkpoints_user_id_slug_key"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt time.Time `gorm:"not null;autoUpdateTime"`
}

func (GormCheckPoint) TableName() string  {
	return "checkpoints"
}

func (g *GormCheckPoint) BeforeCreate(_ *gorm.DB) error {
	if g.ID == uuid.Nil {
		g.ID = uuid.New()
	}
	
	return nil
}
