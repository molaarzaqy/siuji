package converter

import (
	"siuji-backend/internal/entity"
	"siuji-backend/internal/model"
)

func SectionToResponse(section *entity.Section) *model.SectionResponse {
	return &model.SectionResponse{
		PublicID:  section.PublicID.String(),
		Title:     section.Title,
		CreatedAt: section.CreatedAt,
		UpdatedAt: section.UpdatedAt,
	}
}