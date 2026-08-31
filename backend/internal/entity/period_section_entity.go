package entity

import (
	"time"

	"github.com/google/uuid"
)

type PeriodSection struct {
	ID        uint      `gorm:"column:id;primaryKey"`
	PublicID  uuid.UUID `gorm:"column:public_id"`
	PeriodID  uint      `gorm:"column:period_id"`
	SectionID uint      `gorm:"column:section_id"`
	Position  int       `gorm:"column:position"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`

	// Relations
	Period  Period  `gorm:"foreignKey:PeriodID"`
	Section Section `gorm:"foreignKey:SectionID"`
}

func (PeriodSection) TableName() string {
	return "period_sections"
}