package converter

import (
	"siuji-backend/internal/entity"
	"siuji-backend/internal/model"
)

// ParticipantToResponse requires the User relation to be preloaded.
func ParticipantToResponse(pp *entity.ParticipantPeriod) *model.ParticipantResponse {
	return &model.ParticipantResponse{
		PublicID:  pp.PublicID.String(),
		User:      *UserToResponse(&pp.User),
		Status:    pp.Status,
		Score:     pp.Score,
		CreatedAt: pp.CreatedAt,
	}
}

func ParticipantToDetailResponse(pp *entity.ParticipantPeriod, periodPublicID string) *model.ParticipantResponse {
	response := ParticipantToResponse(pp)
	response.PeriodPublicID = periodPublicID
	response.UpdatedAt = pp.UpdatedAt
	return response
}