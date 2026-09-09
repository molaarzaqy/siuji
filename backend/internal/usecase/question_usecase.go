package usecase

import (
	"context"
	"io"
	"siuji-backend/internal/entity"
	"siuji-backend/internal/model"
	"siuji-backend/internal/model/converter"
	"siuji-backend/internal/repository"
	"siuji-backend/pkg/cloudinary"

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
	CloudinaryService *cloudinary.Service
}

func NewQuestionUseCase(
	log *logrus.Logger,
	validate *validator.Validate,
	questionRepository repository.QuestionRepository,
	sectionRepository repository.SectionRepository,
	cloudinaryService *cloudinary.Service,
) *QuestionUseCase {
	return &QuestionUseCase{
		Log:                log,
		Validate:           validate,
		QuestionRepository: questionRepository,
		SectionRepository:  sectionRepository,
		CloudinaryService: cloudinaryService,
	}
}

func (c *QuestionUseCase) Create(ctx context.Context, sectionPublicID string, request *model.QuestionRequest, audioFile, imageFile io.Reader) (*model.QuestionResponse, error) {
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("invalid create question request: %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid request")
	}
	section, err := c.SectionRepository.FindByPublicID(sectionPublicID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "section not found")
	}
	maxNumber, err := c.QuestionRepository.GetMaxNumberInSection(section.ID)
	if err != nil {
		c.Log.Errorf("failed to get max Number: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to create question")
	}
	question := &entity.Question{
		PublicID: uuid.New(),
		SectionID: section.ID,
		Question: request.Question,
		Passage: request.Passage,
		Number: maxNumber + 1,
	}
		if audioFile != nil {
		audioURL, err := c.CloudinaryService.UploadQuestionAudio(ctx, audioFile)
		if err != nil {
			c.Log.Errorf("failed to upload question audio: %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to upload audio file")
		}
		question.AudioURL = &audioURL
	}

	if imageFile != nil {
		imageURL, err := c.CloudinaryService.UploadQuestionImage(ctx, imageFile)
		if err != nil {
			c.Log.Errorf("failed to upload question image: %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to upload image file")
		}
		question.ImageURL = &imageURL
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

func (c *QuestionUseCase) Update(ctx context.Context, publicID string, request *model.QuestionRequest, audioFile, imageFile io.Reader) (*model.QuestionResponse, error) {
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("invalid update question request: %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid request")
	}

	question, err := c.QuestionRepository.FindByPublicID(publicID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "question not found")
	}

	question.Question = request.Question
	question.Passage = request.Passage

	if audioFile != nil {
		audioURL, err := c.CloudinaryService.UploadQuestionAudio(ctx, audioFile)
		if err != nil {
			c.Log.Errorf("failed to upload question audio: %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to upload audio file")
		}
		question.AudioURL = &audioURL
	}

	if imageFile != nil {
		imageURL, err := c.CloudinaryService.UploadQuestionImage(ctx, imageFile)
		if err != nil {
			c.Log.Errorf("failed to upload question image: %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to upload image file")
		}
		question.ImageURL = &imageURL
	}

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

	if err := c.QuestionRepository.UpdateNumbersByPublicIDs(request.QuestionPublicIDs); err != nil {
		c.Log.Errorf("failed to reorder questions: %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, "failed to reorder questions")
	}
	return nil
}