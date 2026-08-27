package models

import "time"

type PeriodSection struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	PublicID  string    `json:"public_id" gorm:"type:varchar(36);uniqueIndex;not null"`
	PeriodID  uint      `json:"period_id" gorm:"not null;index"`
	SectionID uint      `json:"section_id" gorm:"not null;index"`
	Position  int       `json:"position" gorm:"type:int;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relations
	Period  Period  `json:"period,omitempty" gorm:"foreignKey:PeriodID"`
	Section Section `json:"section,omitempty" gorm:"foreignKey:SectionID"`
}

type PeriodSectionResponse struct {
	PeriodSectionPublicID string `json:"period_section_public_id"`
	SectionPublicID       string `json:"section_public_id"`
	Title                 string `json:"title"`
	Position              int    `json:"position"`
}

type ExamStartResponse struct {
	PeriodPublicID string             `json:"period_public_id"`
	Title          string             `json:"title"`
	Sections       []SectionExamItem  `json:"sections"`
}

type SectionExamItem struct {
	PublicID  string              `json:"public_id"`
	Title     string              `json:"title"`
	Position  int                 `json:"position"`
	Questions []QuestionExamItem  `json:"questions"`
}

type QuestionExamItem struct {
	PublicID  string             `json:"public_id"`
	Question  string             `json:"question"`
	AudioURL  *string            `json:"audio_url"`
	ImageURL  *string            `json:"image_url"`
	Passage   *string            `json:"passage"`
	Position  int                `json:"position"`
	Options   []OptionExamItem   `json:"options"`
}

type OptionExamItem struct {
	PublicID   string `json:"public_id"`
	Label      string `json:"label"`
	OptionText string `json:"option_text"`
	Position   int    `json:"position"`
}