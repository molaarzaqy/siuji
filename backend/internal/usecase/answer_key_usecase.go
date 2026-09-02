package usecase

import (
	"siuji-backend/internal/model"
	"siuji-backend/internal/model/converter"
	"siuji-backend/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
)

type AnswerKeyUseCase struct {
	Log                 *logrus.Logger
	Validate            *validator.Validate
	AnswerKeyRepository repository.AnswerKeyRepository
	QuestionRepository  repository.QuestionRepository
	OptionRepository    repository.OptionRepository
}

func NewAnswerKeyUseCase(
	log *logrus.Logger,
	validate *validator.Validate,
	answerKeyRepository repository.AnswerKeyRepository,
	questionRepository repository.QuestionRepository,
	optionRepository repository.OptionRepository,
) *AnswerKeyUseCase {
	return &AnswerKeyUseCase{
		Log:                 log,
		Validate:            validate,
		AnswerKeyRepository: answerKeyRepository,
		QuestionRepository:  questionRepository,
		OptionRepository:    optionRepository,
	}
}

func (c *AnswerKeyUseCase) Upsert(questionPublicID string, request *model.UpsertAnswerKeyRequest) (*model.AnswerKeyResponse, error) {
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("invalid upsert answer key request: %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid request")
	}

	question, err := c.QuestionRepository.FindByPublicID(questionPublicID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "question not found")
	}

	option, err := c.OptionRepository.FindByPublicID(request.OptionPublicID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "option not found")
	}

	if option.QuestionID != question.ID {
		return nil, fiber.NewError(fiber.StatusBadRequest, "option does not belong to this question")
	}

	answerKey, err := c.AnswerKeyRepository.Upsert(question.ID, option.ID)
	if err != nil {
		c.Log.Errorf("failed to upsert answer key: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to save answer key")
	}

	return converter.AnswerKeyToResponse(answerKey, questionPublicID), nil
}