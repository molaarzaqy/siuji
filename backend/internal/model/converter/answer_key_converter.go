package converter

import (
	"siuji-backend/internal/entity"
	"siuji-backend/internal/model"
)

// AnswerKeyToResponse requires the Option relation to be preloaded.
func AnswerKeyToResponse(ak *entity.AnswerKey, questionPublicID string) *model.AnswerKeyResponse {
	return &model.AnswerKeyResponse{
		PublicID:              ak.PublicID.String(),
		QuestionPublicID:      questionPublicID,
		CorrectOptionPublicID: ak.Option.PublicID.String(),
		UpdatedAt:             ak.UpdatedAt,
	}
}