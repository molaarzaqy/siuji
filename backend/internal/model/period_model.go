package model

import "time"

type PeriodRequest struct {
	Title               string    `json:"title" validate:"required"`
	Month               string    `json:"month" validate:"required"`
	Year                int       `json:"year" validate:"required"`
	Status              string    `json:"status" validate:"required,oneof=draft published closed"`
	CertificateURL      string    `json:"certificate_url"`
	CertificateExpMonth time.Time `json:"certificate_exp_month"`
	MinPassingGrade     int       `json:"min_passing_grade"`
	MaxPassingGrade     int       `json:"max_passing_grade"`
	StartTime           time.Time `json:"start_time" validate:"required"`
	EndTime             time.Time `json:"end_time" validate:"required"`
}

type PeriodResponse struct {
	PublicID            string    `json:"public_id"`
	Title               string    `json:"title"`
	Month               string    `json:"month"`
	Year                int       `json:"year"`
	Status              string    `json:"status"`
	CertificateURL      string    `json:"certificate_url,omitempty"`
	CertificateExpMonth time.Time `json:"certificate_exp_month,omitempty"`
	MinPassingGrade     int       `json:"min_passing_grade,omitempty"`
	MaxPassingGrade     int       `json:"max_passing_grade,omitempty"`
	StartTime           time.Time `json:"start_time"`
	EndTime             time.Time `json:"end_time"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type PeriodDetailResponse struct {
	PeriodResponse
	Sections []PeriodSectionResponse `json:"sections"`
}