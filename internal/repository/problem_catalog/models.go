package problemcatalog

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Problem struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	Name       string     `gorm:"type:varchar(255);not null"`
	Slug       string     `gorm:"type:varchar(255);not null;uniqueIndex"`
	Difficulty string     `gorm:"type:varchar(32);not null;index"`
	Tag        string     `gorm:"type:varchar(64);index"`
	Status     string     `gorm:"type:varchar(20);not null;default:'draft';index"`
	AuthorID   uuid.UUID  `gorm:"type:uuid;not null;index"`
	VerifiedAt *time.Time `gorm:"index"`
	CreatedAt  time.Time  `gorm:"not null;autoCreateTime"`
	UpdatedAt  time.Time  `gorm:"not null;autoUpdateTime"`

	Content *ProblemContent `gorm:"foreignKey:ProblemID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

func (Problem) TableName() string {
	return "problems"
}

func (p *Problem) BeforeCreate(_ *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

type ProblemContent struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey"`
	ProblemID      uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	DescriptionMD  string    `gorm:"type:text;not null"`
	InputFormatMD  string    `gorm:"type:text"`
	OutputFormatMD string    `gorm:"type:text"`
	ConstraintsMD  string    `gorm:"type:text"`
	NotesMD        string    `gorm:"type:text"`
	CreatedAt      time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt      time.Time `gorm:"not null;autoUpdateTime"`
	Problem        Problem   `gorm:"foreignKey:ProblemID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

func (ProblemContent) TableName() string {
	return "problem_contents"
}

func (pc *ProblemContent) BeforeCreate(_ *gorm.DB) error {
	if pc.ID == uuid.Nil {
		pc.ID = uuid.New()
	}
	return nil
}

type ProblemFilter struct {
	Tag        string
	Difficulty string
	Limit      int
	Offset     int
}
