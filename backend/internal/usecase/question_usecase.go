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
		uploaded, err := c.CloudinaryService.UploadQuestionAudio(ctx, audioFile)
		if err != nil {
			c.Log.Errorf("failed to upload question audio: %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to upload audio file")
		}
		question.AudioURL = &uploaded.URL
		question.AudioPublicID = &uploaded.PublicID
	}

	if imageFile != nil {
		uploaded, err := c.CloudinaryService.UploadQuestionImage(ctx, imageFile)
		if err != nil {
			c.Log.Errorf("failed to upload question image: %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to upload image file")
		}
		question.ImageURL = &uploaded.URL
		question.ImagePublicID = &uploaded.PublicID
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

	var oldAudioID, oldImageID *string
	if question.AudioPublicID != nil {
		oldAudioID = question.AudioPublicID
	}
	if question.ImagePublicID != nil {
		oldImageID = question.ImagePublicID
	}
	var replacedAudio, replacedImage bool

	question.Question = request.Question
	question.Passage = request.Passage

	if audioFile != nil {
		uploaded, err := c.CloudinaryService.UploadQuestionAudio(ctx, audioFile)
		if err != nil {
			c.Log.Errorf("failed to upload question audio: %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to upload audio file")
		}
		question.AudioURL = &uploaded.URL
		question.AudioPublicID = &uploaded.PublicID
		replacedAudio = true
	}

	if imageFile != nil {
		uploaded, err := c.CloudinaryService.UploadQuestionImage(ctx, imageFile)
		if err != nil {
			c.Log.Errorf("failed to upload question image: %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to upload image file")
		}
		question.ImageURL = &uploaded.URL
		question.ImagePublicID = &uploaded.PublicID
		replacedImage = true
	}

	if err := c.QuestionRepository.Update(question); err != nil {
		c.Log.Errorf("failed to update question: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to update question")
	}
	// Best-effort: file lama sudah tidak direferensikan siapapun
	if replacedAudio && oldAudioID != nil {
		if err := c.CloudinaryService.Destroy(ctx, *oldAudioID, cloudinary.ResourceTypeVideo); err != nil {
			c.Log.Errorf("failed to destroy replaced audio %s: %+v", *oldAudioID, err)
		}
	}
	if replacedImage && oldImageID != nil {
		if err := c.CloudinaryService.Destroy(ctx, *oldImageID, cloudinary.ResourceTypeImage); err != nil {
			c.Log.Errorf("failed to destroy replaced image %s: %+v", *oldImageID, err)
		}
	}

	return converter.QuestionToResponse(question), nil
}


func (c *QuestionUseCase) Delete(ctx context.Context, publicID string) error {
	question, err := c.QuestionRepository.FindByPublicID(publicID)
	if err != nil {
		return mapNotFoundError(c.Log, err, "failed to find question for deletion")
	}

	if err := c.QuestionRepository.Delete(publicID); err != nil {
		return mapNotFoundError(c.Log, err, "failed to delete question")
	}

	if question.AudioPublicID != nil {
		if err := c.CloudinaryService.Destroy(ctx, *question.AudioPublicID, cloudinary.ResourceTypeVideo); err != nil {
			c.Log.Errorf("failed to destroy audio %s: %+v", *question.AudioPublicID, err)
		}
	}
	if question.ImagePublicID != nil {
		if err := c.CloudinaryService.Destroy(ctx, *question.ImagePublicID, cloudinary.ResourceTypeImage); err != nil {
			c.Log.Errorf("failed to destroy image %s: %+v", *question.ImagePublicID, err)
		}
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