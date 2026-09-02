package entity

import (
	"time"

	"github.com/google/uuid"
)

type AnswerKey struct {
	ID         uint      `gorm:"column:id;primaryKey"`
	PublicID   uuid.UUID `gorm:"column:public_id"`
	OptionID   uint      `gorm:"column:option_id"`
	QuestionID uint      `gorm:"column:question_id"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`

	// Relations
	Option Option `gorm:"foreignKey:OptionID"`
}

func (AnswerKey) TableName() string {
	return "answer_keys"
}