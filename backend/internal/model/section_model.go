package model

import "time"

type SectionRequest struct {
	Title string `json:"title" validate:"required"`
}

type SectionResponse struct {
	PublicID  string    `json:"public_id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}