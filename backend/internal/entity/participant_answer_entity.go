package entity

import (
	"time"

	"github.com/google/uuid"
)

type ParticipantAnswer struct {
	ID                  uint              `gorm:"column:id;primaryKey"`
	PublicID            uuid.UUID         `gorm:"column:public_id"`
	ParticipantPeriodID uint              `gorm:"column:participant_period_id"`
	QuestionID          uint              `gorm:"column:question_id"`
	OptionID            uint              `gorm:"column:option_id"`
	IsCorrect           bool              `gorm:"column:is_correct"`
	CreatedAt           time.Time         `gorm:"column:created_at"`
	UpdatedAt           time.Time         `gorm:"column:updated_at"`

	// Relations
	ParticipantPeriod   ParticipantPeriod `gorm:"foreignKey:ParticipantPeriodID"`
	Question            Question          `gorm:"foreignKey:QuestionID"`
	Option              Option            `gorm:"foreignKey:OptionID"`
}

func (ParticipantAnswer) TableName() string {
	return "participant_answers"
}

