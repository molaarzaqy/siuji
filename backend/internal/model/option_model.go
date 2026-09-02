package model

type OptionRequest struct {
	Label      string `json:"label" validate:"required"`
	OptionText string `json:"option_text" validate:"required"`
}

type OptionResponse struct {
	PublicID   string `json:"public_id"`
	Label      string `json:"label"`
	OptionText string `json:"option_text"`
	Position   int    `json:"position"`
}

type ReorderOptionsRequest struct {
	OptionPublicIDs []string `json:"option_public_ids" validate:"required,min=1,dive,uuid"`
}