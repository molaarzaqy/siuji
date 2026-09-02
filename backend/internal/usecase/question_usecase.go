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

type QuestionUseCase struct {
	Log                *logrus.Logger
	Validate           *validator.Validate
	QuestionRepository repository.QuestionRepository
	SectionRepository  repository.SectionRepository
}

func NewQuestionUseCase(
	log *logrus.Logger,
	validate *validator.Validate,
	questionRepository repository.QuestionRepository,
	sectionRepository repository.SectionRepository,
) *QuestionUseCase {
	return &QuestionUseCase{
		Log:                log,
		Validate:           validate,
		QuestionRepository: questionRepository,
		SectionRepository:  sectionRepository,
	}
}

func (c *QuestionUseCase) Create(sectionPublicID string, request *model.QuestionRequest) (*model.QuestionResponse, error) {
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("invalid create question request: %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid request")
	}

	section, err := c.SectionRepository.FindByPublicID(sectionPublicID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "section not found")
	}

	maxPosition, err := c.QuestionRepository.GetMaxPositionInSection(section.ID)
	if err != nil {
		c.Log.Errorf("failed to get max position: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to create question")
	}

	question := &entity.Question{
		PublicID:  uuid.New(),
		SectionID: section.ID,
		Question:  request.Question,
		AudioURL:  request.AudioURL,
		ImageURL:  request.ImageURL,
		Passage:   request.Passage,
		Position:  maxPosition + 1,
	}

	if err := c.QuestionRepository.Create(question); err != nil {
		c.Log.Errorf("failed to create question: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to create question")
	}

	return converter.QuestionToResponse(question), nil
}

func (c *QuestionUseCase) GetDetail(publicID string) (*model.QuestionDetailResponse, error) {
	question, err := c.QuestionRepository.FindByPublicIDWithOptions(publicID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "question not found")
	}
	return converter.QuestionToDetailResponse(question), nil
}

func (c *QuestionUseCase) Update(publicID string, request *model.QuestionRequest) (*model.QuestionResponse, error) {
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("invalid update question request: %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid request")
	}

	question, err := c.QuestionRepository.FindByPublicID(publicID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "question not found")
	}

	question.Question = request.Question
	question.AudioURL = request.AudioURL
	question.ImageURL = request.ImageURL
	question.Passage = request.Passage

	if err := c.QuestionRepository.Update(question); err != nil {
		c.Log.Errorf("failed to update question: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to update question")
	}

	return converter.QuestionToResponse(question), nil
}

func (c *QuestionUseCase) Delete(publicID string) error {
	if err := c.QuestionRepository.Delete(publicID); err != nil {
		return fiber.NewError(fiber.StatusNotFound, "question not found")
	}
	return nil
}

func (c *QuestionUseCase) Reorder(request *model.ReorderQuestionsRequest) error {
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("invalid reorder request: %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, "invalid request")
	}

	if err := c.QuestionRepository.UpdatePositionsByPublicIDs(request.QuestionPublicIDs); err != nil {
		c.Log.Errorf("failed to reorder questions: %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, "failed to reorder questions")
	}
	return nil
}