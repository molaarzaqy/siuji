package converter

import (
	"siuji-backend/internal/entity"
	"siuji-backend/internal/model"
)

func QuestionToResponse(question *entity.Question) *model.QuestionResponse {
	return &model.QuestionResponse{
		PublicID: question.PublicID.String(),
		Question: question.Question,
		AudioURL: question.AudioURL,
		ImageURL: question.ImageURL,
		Passage:  question.Passage,
		Position: question.Position,
	}
}

// QuestionToDetailResponse requires Options and AnswerKeys.Option to be preloaded.
func QuestionToDetailResponse(question *entity.Question) *model.QuestionDetailResponse {
	options := make([]model.OptionResponse, 0, len(question.Options))
	for _, opt := range question.Options {
		options = append(options, *OptionToResponse(&opt))
	}

	var correctOptionPublicID *string
	if len(question.AnswerKeys) > 0 {
		id := question.AnswerKeys[0].Option.PublicID.String()
		correctOptionPublicID = &id
	}

	return &model.QuestionDetailResponse{
		QuestionResponse:       *QuestionToResponse(question),
		CorrectOptionPublicID: correctOptionPublicID,
		Options:                options,
	}
}