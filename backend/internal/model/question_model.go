package model

type QuestionRequest struct {
	Question string  `json:"question" validate:"required"`
	AudioURL *string `json:"audio_url"`
	ImageURL *string `json:"image_url"`
	Passage  *string `json:"passage"`
}

type QuestionResponse struct {
	PublicID string  `json:"public_id"`
	Question string  `json:"question"`
	AudioURL *string `json:"audio_url,omitempty"`
	ImageURL *string `json:"image_url,omitempty"`
	Passage  *string `json:"passage,omitempty"`
	Position int     `json:"position"`
}

type QuestionDetailResponse struct {
	QuestionResponse
	CorrectOptionPublicID *string          `json:"correct_option_public_id"`
	Options               []OptionResponse `json:"options"`
}

type ReorderQuestionsRequest struct {
	QuestionPublicIDs []string `json:"question_public_ids" validate:"required,min=1,dive,uuid"`
}