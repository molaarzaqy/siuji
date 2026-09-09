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

// LabelFromPosition derives an alphabetical option label (A, B, C, ...) from
// a 1-indexed position, so labels always stay consistent with position —
// no more manual labels drifting out of sync after a reorder.
// Wraps to AA, AB, ... beyond Z (unlikely in practice for exam options).
func LabelFromPosition(position int) string {
	if position < 1 {
		position = 1
	}
	label := ""
	for position > 0 {
		position--
		label = string(rune('A'+(position%26))) + label
		position /= 26
	}
	return label
}