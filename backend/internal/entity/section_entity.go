package entity

import (
	"time"

	"github.com/google/uuid"
)

type Section struct {
	ID        uint      `gorm:"column:id;primaryKey"`
	PublicID  uuid.UUID `gorm:"column:public_id"`
	Title     string    `gorm:"column:title"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`

	// Relations
	PeriodSections []PeriodSection `gorm:"foreignKey:SectionID"`
}

func (Section) TableName() string {
	return "sections"
}