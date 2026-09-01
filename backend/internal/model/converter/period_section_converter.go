package converter

import (
	"siuji-backend/internal/entity"
	"siuji-backend/internal/model"
)

// PeriodSectionToResponse requires the Section relation to be preloaded on the entity.
func PeriodSectionToResponse(ps *entity.PeriodSection) *model.PeriodSectionResponse {
	return &model.PeriodSectionResponse{
		PeriodSectionPublicID: ps.PublicID.String(),
		SectionPublicID:       ps.Section.PublicID.String(),
		Title:                 ps.Section.Title,
		Position:              ps.Position,
	}
}