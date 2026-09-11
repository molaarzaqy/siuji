package model

type ScoreConversionRequest struct {
	SectionType  string `json:"section_type" validate:"required,oneof=listening structure reading"`
	CorrectCount int    `json:"correct_count" validate:"min=0"`
	ScaledScore  int    `json:"scaled_score" validate:"min=0"`
}

type ScoreConversionResponse struct {
	ID           uint   `json:"id"`
	SectionType  string `json:"section_type"`
	CorrectCount int    `json:"correct_count"`
	ScaledScore  int    `json:"scaled_score"`
}

type BulkScoreConversionRequest struct {
	Conversions []ScoreConversionRequest `json:"conversions" validate:"required,min=1,dive"`
}

type BulkScoreConversionResponse struct {
	TotalSaved int `json:"total_saved"`
}