package model

type AssignSectionRequest struct {
	SectionPublicID string `json:"section_public_id" validate:"required,uuid"`
	Position        int    `json:"position" validate:"required,min=1"`
}

type ReorderSectionsRequest struct {
	SectionPublicIDs []string `json:"section_public_ids" validate:"required,min=1,dive,uuid"`
}

type PeriodSectionResponse struct {
	PeriodSectionPublicID string `json:"period_section_public_id"`
	SectionPublicID       string `json:"section_public_id"`
	Title                 string `json:"title"`
	Position              int    `json:"position"`
}