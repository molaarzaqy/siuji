package model

import "time"

// Endpoint 1: GET /api/v1/participant/periods
type ParticipantPeriodListResponse struct {
	PublicID string             `json:"public_id"`
	Period   PeriodSimpleInfo   `json:"period"`
	Status   string             `json:"status"`
	Score    *int               `json:"score"`
}

type PeriodSimpleInfo struct {
	PublicID  string    `json:"public_id"`
	Title     string    `json:"title"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Status    string    `json:"status"`
}

// Endpoint 2: GET /api/v1/participant/periods/:period_public_id
type ParticipantPeriodDetailResponse struct {
	PublicID          string    `json:"public_id"`
	Title             string    `json:"title"`
	Month             string    `json:"month"`
	Year              int       `json:"year"`
	StartTime         time.Time `json:"start_time"`
	EndTime           time.Time `json:"end_time"`
	ParticipantStatus string    `json:"participant_status"`
}

// Endpoint 3: POST /api/v1/participant/periods/:period_public_id/start
type ExamSessionResponse struct {
	PeriodPublicID string                `json:"period_public_id"`
	Title          string                `json:"title"`
	Sections       []ExamSectionResponse `json:"sections"`
}

type ExamSectionResponse struct {
	PublicID  string                 `json:"public_id"`
	Title     string                 `json:"title"`
	Position  int                    `json:"position"`
	Questions []ExamQuestionResponse `json:"questions"`
}

type ExamQuestionResponse struct {
	PublicID string               `json:"public_id"`
	Question string               `json:"question"`
	AudioURL *string              `json:"audio_url"`
	ImageURL *string              `json:"image_url"`
	Passage  *string              `json:"passage"`
	Number   int                  `json:"number"`
	Options  []ExamOptionResponse `json:"options"`
}

type ExamOptionResponse struct {
	PublicID   string `json:"public_id"`
	Label      string `json:"label"`
	OptionText string `json:"option_text"`
	Position   int    `json:"position"`
}

// Endpoint 4: POST /api/v1/participant/periods/:period_public_id/answers
type SaveAnswerRequest struct {
	QuestionPublicID string `json:"question_public_id" validate:"required"`
	OptionPublicID   string `json:"option_public_id" validate:"required"`
}

type SaveAnswerResponse struct {
	QuestionPublicID string    `json:"question_public_id"`
	OptionPublicID   string    `json:"option_public_id"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// Endpoint 5: POST /api/v1/participant/periods/:period_public_id/submit
type SubmitExamResponse struct {
	PeriodPublicID string    `json:"period_public_id"`
	Status         string    `json:"status"`
	SubmittedAt    time.Time `json:"submitted_at"`
}

// Endpoint 6: GET /api/v1/participant/periods/:period_public_id/result
type ExamResultResponse struct {
	PeriodTitle       string                    `json:"period_title"`
	Status            string                    `json:"status"`
	FinalScore        int                       `json:"final_score"`
	MinPassingGrade   int                       `json:"min_passing_grade"`
	CertificateURL    *string                   `json:"certificate_url"`
	SectionBreakdowns []SectionBreakdownResponse `json:"section_breakdowns"`
}

type SectionBreakdownResponse struct {
	SectionTitle string `json:"section_title"`
	CorrectCount int    `json:"correct_count"`
	ScaledScore  int    `json:"scaled_score"`
}

