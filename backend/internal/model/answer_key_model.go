package model

import "time"

type UpsertAnswerKeyRequest struct {
	OptionPublicID string `json:"option_public_id" validate:"required,uuid"`
}

type AnswerKeyResponse struct {
	PublicID              string    `json:"public_id"`
	QuestionPublicID      string    `json:"question_public_id"`
	CorrectOptionPublicID string    `json:"correct_option_public_id"`
	UpdatedAt             time.Time `json:"updated_at"`
}