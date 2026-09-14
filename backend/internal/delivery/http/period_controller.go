package http

import (
	"strconv"
	"time"

	"siuji-backend/internal/model"
	"siuji-backend/internal/usecase"
	"siuji-backend/pkg/response"

	"github.com/gofiber/fiber/v3"
)

type PeriodController struct {
	UseCase *usecase.PeriodUseCase
}

func NewPeriodController(useCase *usecase.PeriodUseCase) *PeriodController {
	return &PeriodController{UseCase: useCase}
}

var certificateTemplateTypes = []string{"image/jpeg", "image/png", "application/pdf"}

// Create godoc
// @Summary      Create period
// @Description  Create a new exam period with a certificate template image.
// @Tags         Period
// @Accept       mpfd
// @Produce      json
// @Security     BearerAuth
// @Param        title formData string true "Period title"
// @Param        month formData string true "Month"
// @Param        year formData int true "Year"
// @Param        status formData string true "draft, published, or closed"
// @Param        certificate_exp_month formData string false "RFC3339 datetime"
// @Param        min_passing_grade formData int false "Minimum passing grade"
// @Param        max_passing_grade formData int false "Maximum passing grade"
// @Param        start_time formData string true "RFC3339 datetime"
// @Param        end_time formData string true "RFC3339 datetime"
// @Param        certificate_template formData file true "Certificate template image (JPEG/PNG)"
// @Success      201 {object} response.Response{data=model.PeriodResponse}
// @Failure      400 {object} response.ResponseNoData
// @Router       /periods [post]
func (ctrl *PeriodController) Create(c fiber.Ctx) error {
	year, _ := strconv.Atoi(c.FormValue("year"))
	minPassingGrade, _ := strconv.Atoi(c.FormValue("min_passing_grade"))
	maxPassingGrade, _ := strconv.Atoi(c.FormValue("max_passing_grade"))

	startTime, err := time.Parse(time.RFC3339, c.FormValue("start_time"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid start_time format, use RFC3339 (e.g. 2026-09-01T09:00:00Z)")
	}
	endTime, err := time.Parse(time.RFC3339, c.FormValue("end_time"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid end_time format, use RFC3339")
	}

	var certificateExpMonth time.Time
	if raw := c.FormValue("certificate_exp_month"); raw != "" {
		certificateExpMonth, err = time.Parse(time.RFC3339, raw)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid certificate_exp_month format, use RFC3339")
		}
	}
	request := &model.PeriodRequest{
		Title:               c.FormValue("title"),
		Month:               c.FormValue("month"),
		Year:                year,
		Status:              c.FormValue("status"),
		CertificateExpMonth: certificateExpMonth,
		MinPassingGrade:     minPassingGrade,
		MaxPassingGrade:     maxPassingGrade,
		StartTime:           startTime,
		EndTime:             endTime,
	}
	fileHeader, err := c.FormFile("certificate_template")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "certificate_template file is required")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to read uploaded file")
	}
	defer file.Close()

	validatedFile, err := detectAndValidateContentType(file, certificateTemplateTypes, "certificate_template")
	if err != nil {
		return err
	}

	result, err := ctrl.UseCase.Create(c, request, validatedFile)
	if err != nil {
		return err
	}
	return response.Created(c, "Period created successfully.", result)
}

// GetAll godoc
// @Summary      List periods
// @Description  Get a paginated list of exam periods.
// @Tags         Period
// @Produce      json
// @Security     BearerAuth
// @Param        page query int false "Page number" default(1)
// @Param        limit query int false "Items per page" default(10)
// @Param        filter query string false "Search by title"
// @Param        sort query string false "Sort field, e.g. -created_at"
// @Success      200 {object} response.ResponsePaginated{data=[]model.PeriodResponse}
// @Router       /periods [get]
func (ctrl *PeriodController) GetAll(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	filter := c.Query("filter", "")
	sort := c.Query("sort", "")
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	periods, total, err := ctrl.UseCase.GetAll(filter, sort, limit, offset)
	if err != nil {
		return err
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	meta := response.PaginationMeta{Page: page, Limit: limit, TotalDatas: int(total), TotalPages: totalPages, Filter: filter, Sort: sort}
	return response.SuccessPagination(c, "List period retrieved successfully.", periods, meta)
}

// GetDetail godoc
// @Summary      Get period detail
// @Description  Get a period with its assigned sections.
// @Tags         Period
// @Produce      json
// @Security     BearerAuth
// @Param        period_public_id path string true "Period public ID"
// @Success      200 {object} response.Response{data=model.PeriodDetailResponse}
// @Failure      404 {object} response.ResponseNoData
// @Router       /periods/{period_public_id} [get]
func (ctrl *PeriodController) GetDetail(c fiber.Ctx) error {
	result, err := ctrl.UseCase.GetDetail(c.Params("period_public_id"))
	if err != nil {
		return err
	}
	return response.Success(c, "Period detail retrieved successfully.", result)
}

// Update godoc
// @Summary      Update period
// @Description  Update a period. certificate_template is optional — omit to keep the existing one.
// @Tags         Period
// @Accept       mpfd
// @Produce      json
// @Security     BearerAuth
// @Param        period_public_id path string true "Period public ID"
// @Param        title formData string true "Period title"
// @Param        month formData string true "Month"
// @Param        year formData int true "Year"
// @Param        status formData string true "draft, published, or closed"
// @Param        certificate_exp_month formData string false "RFC3339 datetime"
// @Param        min_passing_grade formData int false "Minimum passing grade"
// @Param        max_passing_grade formData int false "Maximum passing grade"
// @Param        start_time formData string true "RFC3339 datetime"
// @Param        end_time formData string true "RFC3339 datetime"
// @Param        certificate_template formData file false "New certificate template (optional)"
// @Success      200 {object} response.Response{data=model.PeriodResponse}
// @Failure      404 {object} response.ResponseNoData
// @Router       /periods/{period_public_id} [put]
func (ctrl *PeriodController) Update(c fiber.Ctx) error {
	year, _ := strconv.Atoi(c.FormValue("year"))
	minPassingGrade, _ := strconv.Atoi(c.FormValue("min_passing_grade"))
	maxPassingGrade, _ := strconv.Atoi(c.FormValue("max_passing_grade"))

	startTime, err := time.Parse(time.RFC3339, c.FormValue("start_time"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid start_time format, use RFC3339")
	}
	endTime, err := time.Parse(time.RFC3339, c.FormValue("end_time"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid end_time format, use RFC3339")
	}

	var certificateExpMonth time.Time
	if raw := c.FormValue("certificate_exp_month"); raw != "" {
		certificateExpMonth, err = time.Parse(time.RFC3339, raw)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid certificate_exp_month format, use RFC3339")
		}
	}

	request := &model.PeriodRequest{
		Title:               c.FormValue("title"),
		Month:               c.FormValue("month"),
		Year:                year,
		Status:              c.FormValue("status"),
		CertificateExpMonth: certificateExpMonth,
		MinPassingGrade:     minPassingGrade,
		MaxPassingGrade:     maxPassingGrade,
		StartTime:           startTime,
		EndTime:             endTime,
	}

	certificateTemplate, closeTemplate, err := extractOptionalFile(c, "certificate_template", certificateTemplateTypes)
	if err != nil {
		return err
	}
	defer closeTemplate()

	result, err := ctrl.UseCase.Update(c, c.Params("period_public_id"), request, certificateTemplate)
	if err != nil {
		return err
	}
	return response.Success(c, "Period updated successfully.", result)
}

// Delete godoc
// @Summary      Delete period
// @Tags         Period
// @Produce      json
// @Security     BearerAuth
// @Param        period_public_id path string true "Period public ID"
// @Success      200 {object} response.ResponseNoData
// @Failure      404 {object} response.ResponseNoData
// @Router       /periods/{period_public_id} [delete]
func (ctrl *PeriodController) Delete(c fiber.Ctx) error {
	if err := ctrl.UseCase.Delete(c.Params("period_public_id")); err != nil {
		return err
	}
	return response.SuccessNoData(c, "Period deleted successfully.")
}

// AddSection godoc
// @Summary      Assign section to period
// @Tags         Period
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        period_public_id path string true "Period public ID"
// @Param        request body model.AssignSectionRequest true "Section to assign"
// @Success      201 {object} response.Response{data=model.PeriodSectionResponse}
// @Failure      404 {object} response.ResponseNoData
// @Failure      409 {object} response.ResponseNoData
// @Router       /periods/{period_public_id}/sections [post]
func (ctrl *PeriodController) AddSection(c fiber.Ctx) error {
	request := new(model.AssignSectionRequest)
	if err := c.Bind().Body(request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	result, err := ctrl.UseCase.AddSection(c.Params("period_public_id"), request)
	if err != nil {
		return err
	}
	return response.Created(c, "Section assigned to period successfully.", result)
}

// RemoveSection godoc
// @Summary      Remove section from period
// @Tags         Period
// @Produce      json
// @Security     BearerAuth
// @Param        period_public_id path string true "Period public ID"
// @Param        section_public_id path string true "Section public ID"
// @Success      200 {object} response.ResponseNoData
// @Failure      404 {object} response.ResponseNoData
// @Router       /periods/{period_public_id}/sections/{section_public_id} [delete]
func (ctrl *PeriodController) RemoveSection(c fiber.Ctx) error {
	err := ctrl.UseCase.RemoveSection(c.Params("period_public_id"), c.Params("section_public_id"))
	if err != nil {
		return err
	}
	return response.SuccessNoData(c, "Section removed from period successfully.")
}

// ReorderSections godoc
// @Summary      Reorder sections within a period
// @Tags         Period
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        period_public_id path string true "Period public ID"
// @Param        request body model.ReorderSectionsRequest true "Ordered list of section public IDs"
// @Success      200 {object} response.ResponseNoData
// @Failure      404 {object} response.ResponseNoData
// @Router       /periods/{period_public_id}/sections/reorder [put]
func (ctrl *PeriodController) ReorderSections(c fiber.Ctx) error {
	request := new(model.ReorderSectionsRequest)
	if err := c.Bind().Body(request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if err := ctrl.UseCase.ReorderSections(c.Params("period_public_id"), request); err != nil {
		return err
	}
	return response.SuccessNoData(c, "Section positions updated successfully.")
}