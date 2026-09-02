package converter

import (
	"siuji-backend/internal/entity"
	"siuji-backend/internal/model"
)

func OptionToResponse(option *entity.Option) *model.OptionResponse {
	return &model.OptionResponse{
		PublicID:   option.PublicID.String(),
		Label:      option.Label,
		OptionText: option.OptionText,
		Position:   option.Position,
	}
}