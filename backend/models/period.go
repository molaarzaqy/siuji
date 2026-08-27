package models

import "time"

type Period struct {
	ID                  uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	PublicID            string    `json:"public_id" gorm:"type:varchar(36);uniqueIndex;not null"`
	Title               string    `json:"title" gorm:"type:varchar(255);not null"`
	Month               string    `json:"month" gorm:"type:varchar(50)"`
	Year                int       `json:"year" gorm:"type:int"`
	DueDate             time.Time `json:"due_date" gorm:"type:timestamp"`
	Status              string    `json:"status" gorm:"type:varchar(50);not null;default:'draft'"`
	CertificateURL      string    `json:"certificate_url" gorm:"type:text"`
	CertificateExpMonth time.Time `json:"certificate_exp_month" gorm:"type:timestamp"`
	MinPassingGrade     int       `json:"min_passing_grade" gorm:"type:int"`
	MaxPassingGrade     int       `json:"max_passing_grade" gorm:"type:int"`
	StartTime           time.Time `json:"start_time" gorm:"type:timestamp"`
	EndTime             time.Time `json:"end_time" gorm:"type:timestamp"`
	CreatedAt           time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt           time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relations
	ParticipantPeriods []ParticipantPeriod `json:"participant_periods,omitempty" gorm:"foreignKey:PeriodID"`
	PeriodSections     []PeriodSection     `json:"period_sections,omitempty" gorm:"foreignKey:PeriodID"`
}

type PeriodRequest struct {
	Title               string    `json:"title" validate:"required"`
	Month               string    `json:"month" validate:"required"`
	Year                int       `json:"year" validate:"required"`
	DueDate             time.Time `json:"due_date" validate:"required"`
	Status              string    `json:"status" validate:"required,oneof=draft published closed"`
	CertificateURL      string    `json:"certificate_url"`
	CertificateExpMonth time.Time `json:"certificate_exp_month"`
	MinPassingGrade     int       `json:"min_passing_grade"`
	MaxPassingGrade     int       `json:"max_passing_grade"`
	StartTime           time.Time `json:"start_time" validate:"required"`
	EndTime             time.Time `json:"end_time" validate:"required"`
}

type PeriodDetailResponse struct {
	PublicID        	string                  `json:"public_id"`
	Title           	string                  `json:"title"`
	Month           	string                  `json:"month"`
	Year            	int                     `json:"year"`
	DueDate         	time.Time               `json:"due_date"`
	Status          	string                  `json:"status"`
	CertificateURL  	string                  `json:"certificate_url"`
	CertificateExpMonth time.Time               `json:"certificate_exp"`
	MinPassingGrade 	int                     `json:"min_passing_grade"`
	MaxPassingGrade 	int                     `json:"max_passing_grade"`
	StartTime       	time.Time               `json:"start_time"`
	EndTime         	time.Time               `json:"end_time"`
	Sections        	[]PeriodSectionResponse `json:"sections"`
	CreatedAt       	time.Time              	`json:"created_at"`
	UpdatedAt       	time.Time              	`json:"updated_at"`
}
