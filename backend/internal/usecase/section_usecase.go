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

type SectionUseCase struct {
	Log               *logrus.Logger
	Validate          *validator.Validate
	SectionRepository repository.SectionRepository
}

func NewSectionUseCase(
	log *logrus.Logger,
	validate *validator.Validate,
	sectionRepository repository.SectionRepository,
) *SectionUseCase {
	return &SectionUseCase{
		Log:               log,
		Validate:          validate,
		SectionRepository: sectionRepository,
	}
}

func (c *SectionUseCase) Create(request *model.SectionRequest) (*model.SectionResponse, error) {
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("invalid create section request: %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid request")
	}

	section := &entity.Section{
		PublicID: uuid.New(),
		Title:    request.Title,
	}

	if err := c.SectionRepository.Create(section); err != nil {
		c.Log.Errorf("failed to create section: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to create section")
	}

	return converter.SectionToResponse(section), nil
}

func (c *SectionUseCase) GetAll() ([]model.SectionResponse, error) {
	sections, err := c.SectionRepository.FindAll()
	if err != nil {
		c.Log.Errorf("failed to fetch sections: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to fetch sections")
	}

	responses := make([]model.SectionResponse, 0, len(sections))
	for _, s := range sections {
		responses = append(responses, *converter.SectionToResponse(&s))
	}
	return responses, nil
}

func (c *SectionUseCase) GetByPublicID(publicID string) (*model.SectionResponse, error) {
	section, err := c.SectionRepository.FindByPublicID(publicID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "section not found")
	}
	return converter.SectionToResponse(section), nil
}

func (c *SectionUseCase) Update(publicID string, request *model.SectionRequest) (*model.SectionResponse, error) {
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("invalid update section request: %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid request")
	}

	section, err := c.SectionRepository.FindByPublicID(publicID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "section not found")
	}

	section.Title = request.Title

	if err := c.SectionRepository.Update(section); err != nil {
		c.Log.Errorf("failed to update section: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to update section")
	}

	return converter.SectionToResponse(section), nil
}

func (c *SectionUseCase) Delete(publicID string) error {
	if err := c.SectionRepository.Delete(publicID); err != nil {
		return fiber.NewError(fiber.StatusNotFound, "section not found")
	}
	return nil
}