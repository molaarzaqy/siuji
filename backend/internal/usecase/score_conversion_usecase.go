package usecase

import (
	"siuji-backend/internal/entity"
	"siuji-backend/internal/model"
	"siuji-backend/internal/model/converter"
	"siuji-backend/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
)

type ScoreConversionUseCase struct {
	Log                       *logrus.Logger
	Validate                  *validator.Validate
	ScoreConversionRepository repository.ScoreConversionRepository
}

func NewScoreConversionUseCase(
	log *logrus.Logger,
	validate *validator.Validate,
	scoreConversionRepository repository.ScoreConversionRepository,
) *ScoreConversionUseCase {
	return &ScoreConversionUseCase{
		Log:                       log,
		Validate:                  validate,
		ScoreConversionRepository: scoreConversionRepository,
	}
}

func (c *ScoreConversionUseCase) Create(request *model.ScoreConversionRequest) (*model.ScoreConversionResponse, error) {
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("invalid create score conversion request: %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid request")
	}

	sc := &entity.ScoreConversion{
		SectionType:  request.SectionType,
		CorrectCount: request.CorrectCount,
		ScaledScore:  request.ScaledScore,
	}

	if err := c.ScoreConversionRepository.Create(sc); err != nil {
		c.Log.Errorf("failed to create score conversion: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to create score conversion (possibly duplicate section_type + correct_count)")
	}

	return converter.ScoreConversionToResponse(sc), nil
}

func (c *ScoreConversionUseCase) GetAll(sectionType string) ([]model.ScoreConversionResponse, error) {
	list, err := c.ScoreConversionRepository.FindAll(sectionType)
	if err != nil {
		c.Log.Errorf("failed to fetch score conversions: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to fetch score conversions")
	}

	responses := make([]model.ScoreConversionResponse, 0, len(list))
	for _, sc := range list {
		responses = append(responses, *converter.ScoreConversionToResponse(&sc))
	}
	return responses, nil
}

func (c *ScoreConversionUseCase) Update(id uint, request *model.ScoreConversionRequest) (*model.ScoreConversionResponse, error) {
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("invalid update score conversion request: %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid request")
	}

	sc, err := c.ScoreConversionRepository.FindByID(id)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "score conversion not found")
	}

	sc.SectionType = request.SectionType
	sc.CorrectCount = request.CorrectCount
	sc.ScaledScore = request.ScaledScore

	if err := c.ScoreConversionRepository.Update(sc); err != nil {
		c.Log.Errorf("failed to update score conversion: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to update score conversion")
	}

	return converter.ScoreConversionToResponse(sc), nil
}

func (c *ScoreConversionUseCase) Delete(id uint) error {
	if err := c.ScoreConversionRepository.Delete(id); err != nil {
		return fiber.NewError(fiber.StatusNotFound, "score conversion not found")
	}
	return nil
}

func (c *ScoreConversionUseCase) BulkCreate(request *model.BulkScoreConversionRequest) (*model.BulkScoreConversionResponse, error) {
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("invalid bulk score conversion request: %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid request: "+err.Error())
	}

	list := make([]entity.ScoreConversion, 0, len(request.Conversions))
	for _, item := range request.Conversions {
		list = append(list, entity.ScoreConversion{
			SectionType:  item.SectionType,
			CorrectCount: item.CorrectCount,
			ScaledScore:  item.ScaledScore,
		})
	}

	if err := c.ScoreConversionRepository.BulkUpsert(list); err != nil {
		c.Log.Errorf("failed to bulk upsert score conversions: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to save score conversions")
	}

	return &model.BulkScoreConversionResponse{TotalSaved: len(list)}, nil
}