package model

import "time"

type SectionRequest struct {
	Title string `json:"title" validate:"required"`
	SectionType string `json:"section_type" validate:"required,oneof=listening structure reading"`
}

type SectionResponse struct {
	PublicID  	string    `json:"public_id"`
	Title     	string    `json:"title"`
	SectionType string	  `json:"section_type"`
	CreatedAt 	time.Time `json:"created_at"`
	UpdatedAt 	time.Time `json:"updated_at"`
}

type SectionDetailResponse struct {
	SectionResponse
	Questions []QuestionDetailResponse `json:"questions"`
}