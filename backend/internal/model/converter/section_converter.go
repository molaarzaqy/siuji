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

func SectionToDetailResponse(section *entity.Section) *model.SectionDetailResponse {
	questions := make([]model.QuestionDetailResponse, 0, len(section.Questions))
	for _, q := range section.Questions {
		questions = append(questions, *QuestionToDetailResponse(&q))
	}
	return &model.SectionDetailResponse{
		SectionResponse: *SectionToResponse(section),
		Questions: questions,
	}
}