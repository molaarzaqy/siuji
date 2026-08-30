package model

import "time"

type UserResponse struct {
	PublicID   string    `json:"public_id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Role       string    `json:"role"`
	NIM        string    `json:"nim,omitempty"`
	University string    `json:"university,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}