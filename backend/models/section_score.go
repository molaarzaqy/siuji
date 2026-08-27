package models

import "time"

type SectionScore struct {
	ID                  uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	PublicID            string    `json:"public_id" gorm:"type:varchar(36);uniqueIndex;not null"`
	ParticipantPeriodID uint      `json:"participant_period_id" gorm:"not null;index"`
	SectionID           uint      `json:"section_id" gorm:"not null;index"`
	CorrectCount        int       `json:"correct_count" gorm:"type:int;default:0"`
	RawScore            int       `json:"raw_score" gorm:"type:int;default:0"`
	ScaledScore         int       `json:"scaled_score" gorm:"type:int;default:0"`
	CreatedAt           time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt           time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relations
	ParticipantPeriod ParticipantPeriod `json:"participant_period,omitempty" gorm:"foreignKey:ParticipantPeriodID"`
	Section           Section           `json:"section,omitempty" gorm:"foreignKey:SectionID"`
}