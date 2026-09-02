package entity

import (
	"time"

	"github.com/google/uuid"
)

type Option struct {
	ID         uint      `gorm:"column:id;primaryKey"`
	PublicID   uuid.UUID `gorm:"column:public_id"`
	QuestionID uint      `gorm:"column:question_id"`
	Label      string    `gorm:"column:label"`
	OptionText string    `gorm:"column:option_text"`
	Position   int       `gorm:"column:position"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}

func (Option) TableName() string {
	return "options"
}