package models

import "time"

type ParticipantPeriod struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	PublicID  string    `json:"public_id" gorm:"type:varchar(36);uniqueIndex;not null"`
	UserID    uint      `json:"user_id" gorm:"not null;index"`
	PeriodID  uint      `json:"period_id" gorm:"not null;index"`
	Status    string    `json:"status" gorm:"type:varchar(50);not null;default:'registered'"`
	Score     *int      `json:"score" gorm:"type:int"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relations
	User               User                `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Period             Period              `json:"period,omitempty" gorm:"foreignKey:PeriodID"`
	ParticipantAnswers []ParticipantAnswer `json:"participant_answers,omitempty" gorm:"foreignKey:ParticipantPeriodID"`
	SectionScores      []SectionScore      `json:"section_scores,omitempty" gorm:"foreignKey:ParticipantPeriodID"`
}

type ParticipantRequest struct {
	Name       string `json:"name" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	NIM        string `json:"nim" validate:"required"`
	University string `json:"university" validate:"required"`
}

type UpdateParticipantRequest struct {
	Name       string   `json:"name" validate:"required"`
	NIM        string   `json:"nim" validate:"required"`
	University string   `json:"university" validate:"required"`
	Status     string   `json:"status" validate:"required"`
	Score      *int     `json:"score"`
}

type ImportErrorDetail struct {
	Row     int    `json:"row"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

type ImportParticipantResponse struct {
	TotalImported int                 `json:"total_imported"`
	TotalSkipped  int                 `json:"total_skipped"`
	Errors        []ImportErrorDetail `json:"errors"`
}

type ParticipantListResponse struct {
	PublicID  string              `json:"public_id"`
	User      UserResponse        `json:"user"`
	Status    string              `json:"status"`
	Score     *int                `json:"score"`
	CreatedAt time.Time           `json:"created_at"`
}

type ParticipantResponse struct {
	PublicID       string       `json:"public_id"`
	PeriodPublicID string       `json:"period_public_id"`
	User           UserResponse `json:"user"`
	Status         string       `json:"status"`
	Score          *int         `json:"score"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

type ParticipantPeriodListResponse struct {
	PublicID string              `json:"public_id"`
	Period   PeriodItemResponse  `json:"period"`
	Status   string              `json:"status"`
	Score    *int            `json:"score"`
}

type PeriodItemResponse struct {
	PublicID  string    `json:"public_id"`
	Title     string    `json:"title"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Status    string    `json:"status"`
}

// DTO untuk detail single period participant
type ParticipantPeriodDetailResponse struct {
	PublicID          string    `json:"public_id"`
	Title             string    `json:"title"`
	Month             string    `json:"month"`
	Year              int       `json:"year"`
	StartTime         time.Time `json:"start_time"`
	EndTime           time.Time `json:"end_time"`
	ParticipantStatus string    `json:"participant_status"`
}