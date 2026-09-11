package converter

import (
	"siuji-backend/internal/entity"
	"siuji-backend/internal/model"
)

func ParticipantPeriodToListResponse(pp *entity.ParticipantPeriod) model.ParticipantPeriodListResponse {
	return model.ParticipantPeriodListResponse{
		PublicID: pp.PublicID.String(),
		Period: model.PeriodSimpleInfo{
			PublicID:  pp.Period.PublicID.String(),
			Title:     pp.Period.Title,
			StartTime: pp.Period.StartTime,
			EndTime:   pp.Period.EndTime,
			Status:    pp.Period.Status,
		},
		Status: pp.Status,
		Score:  pp.Score,
	}
}

func PeriodToParticipantDetailResponse(period *entity.Period, participantStatus string) *model.ParticipantPeriodDetailResponse {
	return &model.ParticipantPeriodDetailResponse{
		PublicID:          period.PublicID.String(),
		Title:             period.Title,
		Month:             period.Month,
		Year:              period.Year,
		StartTime:         period.StartTime,
		EndTime:           period.EndTime,
		ParticipantStatus: participantStatus,
	}
}

func QuestionToExamQuestionResponse(q *entity.Question) model.ExamQuestionResponse {
	options := make([]model.ExamOptionResponse, 0, len(q.Options))
	for _, opt := range q.Options {
		options = append(options, model.ExamOptionResponse{
			PublicID:   opt.PublicID.String(),
			Label:      opt.Label,
			OptionText: opt.OptionText,
			Position:   opt.Position,
		})
	}

	return model.ExamQuestionResponse{
		PublicID: q.PublicID.String(),
		Question: q.Question,
		AudioURL: q.AudioURL,
		ImageURL: q.ImageURL,
		Passage:  q.Passage,
		Number:   q.Number,
		Options:  options,
	}
}

