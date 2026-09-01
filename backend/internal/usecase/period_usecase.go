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

type PeriodUseCase struct {
	Log                     *logrus.Logger
	Validate                *validator.Validate
	PeriodRepository        repository.PeriodRepository
	SectionRepository       repository.SectionRepository
	PeriodSectionRepository repository.PeriodSectionRepository
}

func NewPeriodUseCase(
	log *logrus.Logger,
	validate *validator.Validate,
	periodRepository repository.PeriodRepository,
	sectionRepository repository.SectionRepository,
	periodSectionRepository repository.PeriodSectionRepository,
) *PeriodUseCase {
	return &PeriodUseCase{
		Log:                     log,
		Validate:                validate,
		PeriodRepository:        periodRepository,
		SectionRepository:       sectionRepository,
		PeriodSectionRepository: periodSectionRepository,
	}
}

func (c *PeriodUseCase) Create(request *model.PeriodRequest) (*model.PeriodResponse, error) {
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("invalid create period request: %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid request")
	}

	period := &entity.Period{
		PublicID:            uuid.New(),
		Title:               request.Title,
		Month:               request.Month,
		Year:                request.Year,
		Status:              request.Status,
		CertificateURL:      request.CertificateURL,
		CertificateExpMonth: request.CertificateExpMonth,
		MinPassingGrade:     request.MinPassingGrade,
		MaxPassingGrade:     request.MaxPassingGrade,
		StartTime:           request.StartTime,
		EndTime:             request.EndTime,
	}

	if err := c.PeriodRepository.Create(period); err != nil {
		c.Log.Errorf("failed to create period: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to create period")
	}

	return converter.PeriodToResponse(period), nil
}

func (c *PeriodUseCase) GetAll(filter, sort string, limit, offset int) ([]model.PeriodResponse, int64, error) {
	periods, total, err := c.PeriodRepository.FindAllPagination(filter, sort, limit, offset)
	if err != nil {
		c.Log.Errorf("failed to fetch periods: %+v", err)
		return nil, 0, fiber.NewError(fiber.StatusInternalServerError, "failed to fetch periods")
	}

	responses := make([]model.PeriodResponse, 0, len(periods))
	for _, p := range periods {
		responses = append(responses, *converter.PeriodToResponse(&p))
	}
	return responses, total, nil
}

func (c *PeriodUseCase) GetDetail(publicID string) (*model.PeriodDetailResponse, error) {
	period, err := c.PeriodRepository.FindByPublicIDWithSections(publicID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "period not found")
	}
	return converter.PeriodToDetailResponse(period), nil
}

func (c *PeriodUseCase) Update(publicID string, request *model.PeriodRequest) (*model.PeriodResponse, error) {
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("invalid update period request: %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid request")
	}

	period, err := c.PeriodRepository.FindByPublicID(publicID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "period not found")
	}

	period.Title = request.Title
	period.Month = request.Month
	period.Year = request.Year
	period.Status = request.Status
	period.CertificateURL = request.CertificateURL
	period.CertificateExpMonth = request.CertificateExpMonth
	period.MinPassingGrade = request.MinPassingGrade
	period.MaxPassingGrade = request.MaxPassingGrade
	period.StartTime = request.StartTime
	period.EndTime = request.EndTime

	if err := c.PeriodRepository.Update(period); err != nil {
		c.Log.Errorf("failed to update period: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to update period")
	}

	return converter.PeriodToResponse(period), nil
}

func (c *PeriodUseCase) Delete(publicID string) error {
	if err := c.PeriodRepository.Delete(publicID); err != nil {
		return fiber.NewError(fiber.StatusNotFound, "period not found")
	}
	return nil
}

func (c *PeriodUseCase) AddSection(periodPublicID string, request *model.AssignSectionRequest) (*model.PeriodSectionResponse, error) {
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("invalid add section request: %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid request")
	}

	period, err := c.PeriodRepository.FindByPublicID(periodPublicID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "period not found")
	}

	section, err := c.SectionRepository.FindByPublicID(request.SectionPublicID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "section not found")
	}

	exists, err := c.PeriodSectionRepository.ExistsByPeriodAndSection(period.ID, section.ID)
	if err != nil {
		c.Log.Errorf("failed to check existing period section: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to add section")
	}
	if exists {
		return nil, fiber.NewError(fiber.StatusConflict, "section already assigned to this period")
	}

	position := request.Position
	if position <= 0 {
		count, err := c.PeriodSectionRepository.CountByPeriodID(period.ID)
		if err != nil {
			c.Log.Errorf("failed to count period sections: %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to add section")
		}
		position = int(count) + 1
	}

	periodSection := &entity.PeriodSection{
		PublicID:  uuid.New(),
		PeriodID:  period.ID,
		SectionID: section.ID,
		Position:  position,
	}

	if err := c.PeriodSectionRepository.Create(periodSection); err != nil {
		c.Log.Errorf("failed to add section to period: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to add section")
	}

	periodSection.Section = *section // supaya converter dapat Section.Title tanpa query ulang

	return converter.PeriodSectionToResponse(periodSection), nil
}

func (c *PeriodUseCase) RemoveSection(periodPublicID, sectionPublicID string) error {
	period, err := c.PeriodRepository.FindByPublicID(periodPublicID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "period not found")
	}
	section, err := c.SectionRepository.FindByPublicID(sectionPublicID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "section not found")
	}
	if err := c.PeriodSectionRepository.DeleteByPeriodAndSection(period.ID, section.ID); err != nil {
		return fiber.NewError(fiber.StatusNotFound, "section not found in this period")
	}
	return nil
}

func (c *PeriodUseCase) ReorderSections(periodPublicID string, request *model.ReorderSectionsRequest) error {
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("invalid reorder request: %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, "invalid request")
	}

	period, err := c.PeriodRepository.FindByPublicID(periodPublicID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "period not found")
	}

	positions := make(map[uint]int, len(request.SectionPublicIDs))
	for index, sectionPublicID := range request.SectionPublicIDs {
		section, err := c.SectionRepository.FindByPublicID(sectionPublicID)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, "one or more sections not found")
		}
		positions[section.ID] = index + 1
	}

	if err := c.PeriodSectionRepository.UpdatePositionsBulk(period.ID, positions); err != nil {
		c.Log.Errorf("failed to reorder sections: %+v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "failed to reorder sections")
	}

	return nil
}