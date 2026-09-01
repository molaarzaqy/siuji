package model

import "time"

type AddParticipantRequest struct {
	Name       string `json:"name" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	NIM        string `json:"nim" validate:"required"`
	University string `json:"university"`
}

type UpdateParticipantRequest struct {
	Name       string `json:"name"`
	NIM        string `json:"nim"`
	University string `json:"university"`
	Status     string `json:"status" validate:"omitempty,oneof=registered started completed"`
	Score      *int   `json:"score"`
}

type ParticipantResponse struct {
	PublicID       string       `json:"public_id"`
	PeriodPublicID string       `json:"period_public_id,omitempty"`
	User           UserResponse `json:"user"`
	Status         string       `json:"status"`
	Score          *int         `json:"score"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at,omitempty"`
}