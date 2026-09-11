package converter

import (
	"siuji-backend/internal/entity"
	"siuji-backend/internal/model"
)

func ScoreConversionToResponse(sc *entity.ScoreConversion) *model.ScoreConversionResponse {
	return &model.ScoreConversionResponse{
		ID: sc.ID,
		SectionType: sc.SectionType,
		CorrectCount: sc.CorrectCount,
		ScaledScore: sc.ScaledScore,
	}
}