package usecase

import (
	"siuji-backend/internal/entity"
	"siuji-backend/internal/model"
	"siuji-backend/internal/model/converter"
	"siuji-backend/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type OptionUseCase struct {
	Log                *logrus.Logger
	Validate           *validator.Validate
	OptionRepository   repository.OptionRepository
	QuestionRepository repository.QuestionRepository
}

func NewOptionUseCase(
	log *logrus.Logger,
	validate *validator.Validate,
	optionRepository repository.OptionRepository,
	questionRepository repository.QuestionRepository,
) *OptionUseCase {
	return &OptionUseCase{
		Log:                log,
		Validate:           validate,
		OptionRepository:   optionRepository,
		QuestionRepository: questionRepository,
	}
}

func (c *OptionUseCase) Create(questionPublicID string, request *model.OptionRequest) (*model.OptionResponse, error) {
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("invalid create option request: %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid request")
	}

	question, err := c.QuestionRepository.FindByPublicID(questionPublicID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "question not found")
	}

	maxPosition, err := c.OptionRepository.GetMaxPositionInQuestion(question.ID)
	if err != nil {
		c.Log.Errorf("failed to get max position: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to create option")
	}

	option := &entity.Option{
		PublicID:   uuid.New(),
		QuestionID: question.ID,
		Label:      request.Label,
		OptionText: request.OptionText,
		Position:   maxPosition + 1,
	}

	if err := c.OptionRepository.Create(option); err != nil {
		c.Log.Errorf("failed to create option: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to create option")
	}

	return converter.OptionToResponse(option), nil
}

func (c *OptionUseCase) Update(publicID string, request *model.OptionRequest) (*model.OptionResponse, error) {
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("invalid update option request: %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid request")
	}

	option, err := c.OptionRepository.FindByPublicID(publicID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "option not found")
	}

	option.Label = request.Label
	option.OptionText = request.OptionText

	if err := c.OptionRepository.Update(option); err != nil {
		c.Log.Errorf("failed to update option: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to update option")
	}

	return converter.OptionToResponse(option), nil
}

func (c *OptionUseCase) Delete(publicID string) error {
	if err := c.OptionRepository.Delete(publicID); err != nil {
		return fiber.NewError(fiber.StatusNotFound, "option not found")
	}
	return nil
}

func (c *OptionUseCase) Reorder(request *model.ReorderOptionsRequest) error {
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("invalid reorder request: %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, "invalid request")
	}

	if err := c.OptionRepository.UpdatePositionsByPublicIDs(request.OptionPublicIDs); err != nil {
		c.Log.Errorf("failed to reorder options: %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, "failed to reorder options")
	}
	return nil
}