package entity

import (
	"time"

	"github.com/google/uuid"
)

type SectionScore struct {
	ID                  uint              `gorm:"column:id;primaryKey"`
	PublicID            uuid.UUID         `gorm:"column:public_id"`
	ParticipantPeriodID uint              `gorm:"column:participant_period_id"`
	SectionID           uint              `gorm:"column:section_id"`
	CorrectCount        int               `gorm:"column:correct_count"`
	RawScore            int               `gorm:"column:raw_score"`
	ScaledScore         int               `gorm:"column:scaled_score"`
	CreatedAt           time.Time         `gorm:"column:created_at"`
	UpdatedAt           time.Time         `gorm:"column:updated_at"`

	// Relations
	ParticipantPeriod   ParticipantPeriod `gorm:"foreignKey:ParticipantPeriodID"`
	Section             Section           `gorm:"foreignKey:SectionID"`
}

func (SectionScore) TableName() string {
	return "section_scores"
}

