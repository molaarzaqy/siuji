package entity

import (
	"time"

	"github.com/google/uuid"
)

type Question struct {
	ID        uint      `gorm:"column:id;primaryKey"`
	PublicID  uuid.UUID `gorm:"column:public_id"`
	SectionID uint      `gorm:"column:section_id"`
	Question  string    `gorm:"column:question"`
	AudioURL  *string   `gorm:"column:audio_url"`
	ImageURL  *string   `gorm:"column:image_url"`
	Passage   *string   `gorm:"column:passage"`
	Position  int       `gorm:"column:position"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`

	// Relations
	Options    []Option    `gorm:"foreignKey:QuestionID"`
	AnswerKeys []AnswerKey `gorm:"foreignKey:QuestionID"`
}

func (Question) TableName() string {
	return "questions"
}