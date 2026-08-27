package models

import "time"

type ParticipantAnswer struct {
	ID                  uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	PublicID            string    `json:"public_id" gorm:"type:varchar(36);uniqueIndex;not null"`
	ParticipantPeriodID uint      `json:"participant_period_id" gorm:"not null;index"`
	QuestionID          uint      `json:"question_id" gorm:"not null;index"`
	OptionID            uint      `json:"option_id" gorm:"not null;index"`
	IsCorrect           bool      `json:"is_correct" gorm:"type:boolean;default:false"`
	CreatedAt           time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt           time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relations
	ParticipantPeriod ParticipantPeriod `json:"participant_period,omitempty" gorm:"foreignKey:ParticipantPeriodID"`
	Question          Question          `json:"question,omitempty" gorm:"foreignKey:QuestionID"`
	Option            Option            `json:"option,omitempty" gorm:"foreignKey:OptionID"`
}

type SaveAnswerRequest struct {
	QuestionPublicID string `json:"question_public_id" validate:"required"`
	OptionPublicID   string `json:"option_public_id" validate:"required"`
}

// DTO untuk Response Body setelah jawaban berhasil disimpan
type AnswerResponse struct {
	QuestionPublicID string    `json:"question_public_id"`
	OptionPublicID   string    `json:"option_public_id"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// DTO untuk Response Body saat peserta menekan tombol Submit Ujian
type ExamSubmitResponse struct {
	PeriodPublicID string    `json:"period_public_id"`
	Status         string    `json:"status"` // "completed"
	SubmittedAt    time.Time `json:"submitted_at"`
}