package http

import (
	"strconv"

	"siuji-backend/internal/model"
	"siuji-backend/internal/usecase"
	"siuji-backend/pkg/response"

	"github.com/gofiber/fiber/v3"
)

type ParticipantController struct {
	UseCase *usecase.ParticipantUseCase
}

func NewParticipantController(useCase *usecase.ParticipantUseCase) *ParticipantController {
	return &ParticipantController{UseCase: useCase}
}

// Add godoc
// @Summary      Add participant to period
// @Description  Creates the user account if the email doesn't exist yet. Password is auto-generated from NIM.
// @Tags         Participant
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        period_public_id path string true "Period public ID"
// @Param        request body model.AddParticipantRequest true "Participant data"
// @Success      201 {object} response.Response{data=model.ParticipantResponse}
// @Failure      404 {object} response.ResponseNoData
// @Failure      409 {object} response.ResponseNoData
// @Router       /periods/{period_public_id}/participants [post]
func (ctrl *ParticipantController) Add(c fiber.Ctx) error {
	request := new(model.AddParticipantRequest)
	if err := c.Bind().Body(request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	result, err := ctrl.UseCase.AddParticipant(c.Params("period_public_id"), request)
	if err != nil {
		return err
	}
	return response.Created(c, "Participant added to period successfully.", result)
}

// GetAll godoc
// @Summary      List participants in a period
// @Tags         Participant
// @Produce      json
// @Security     BearerAuth
// @Param        period_public_id path string true "Period public ID"
// @Param        page query int false "Page number" default(1)
// @Param        limit query int false "Items per page" default(10)
// @Param        filter query string false "Search by name/email/nim"
// @Param        sort query string false "Sort field"
// @Success      200 {object} response.ResponsePaginated{data=[]model.ParticipantResponse}
// @Failure      404 {object} response.ResponseNoData
// @Router       /periods/{period_public_id}/participants [get]
func (ctrl *ParticipantController) GetAll(c fiber.Ctx) error {
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

	participants, total, err := ctrl.UseCase.GetAllByPeriod(c.Params("period_public_id"), filter, sort, limit, offset)
	if err != nil {
		return err
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	meta := response.PaginationMeta{Page: page, Limit: limit, TotalDatas: int(total), TotalPages: totalPages, Filter: filter, Sort: sort}
	return response.SuccessPagination(c, "List participant retrieved successfully.", participants, meta)
}

// GetDetail godoc
// @Summary      Get participant detail
// @Tags         Participant
// @Produce      json
// @Security     BearerAuth
// @Param        period_public_id path string true "Period public ID"
// @Param        user_public_id path string true "User public ID"
// @Success      200 {object} response.Response{data=model.ParticipantResponse}
// @Failure      404 {object} response.ResponseNoData
// @Router       /periods/{period_public_id}/participants/{user_public_id} [get]
func (ctrl *ParticipantController) GetDetail(c fiber.Ctx) error {
	result, err := ctrl.UseCase.GetDetail(c.Params("period_public_id"), c.Params("user_public_id"))
	if err != nil {
		return err
	}
	return response.Success(c, "Participant detail retrieved successfully.", result)
}

// Update godoc
// @Summary      Update participant
// @Description  Partial update — only non-empty fields are applied. Changing NIM regenerates the password.
// @Tags         Participant
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        period_public_id path string true "Period public ID"
// @Param        user_public_id path string true "User public ID"
// @Param        request body model.UpdateParticipantRequest true "Fields to update"
// @Success      200 {object} response.Response{data=model.ParticipantResponse}
// @Failure      404 {object} response.ResponseNoData
// @Router       /periods/{period_public_id}/participants/{user_public_id} [put]
func (ctrl *ParticipantController) Update(c fiber.Ctx) error {
	request := new(model.UpdateParticipantRequest)
	if err := c.Bind().Body(request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	result, err := ctrl.UseCase.Update(c.Params("period_public_id"), c.Params("user_public_id"), request)
	if err != nil {
		return err
	}
	return response.Success(c, "Participant updated successfully.", result)
}

// Remove godoc
// @Summary      Remove participant from period
// @Tags         Participant
// @Produce      json
// @Security     BearerAuth
// @Param        period_public_id path string true "Period public ID"
// @Param        user_public_id path string true "User public ID"
// @Success      200 {object} response.ResponseNoData
// @Failure      404 {object} response.ResponseNoData
// @Router       /periods/{period_public_id}/participants/{user_public_id} [delete]
func (ctrl *ParticipantController) Remove(c fiber.Ctx) error {
	err := ctrl.UseCase.Remove(c.Params("period_public_id"), c.Params("user_public_id"))
	if err != nil {
		return err
	}
	return response.SuccessNoData(c, "Participant removed from period successfully.")
}

// Import godoc
// @Summary      Import participants from Excel
// @Description  Columns (in order, header row skipped): Name, Email, NIM, University.
// @Tags         Participant
// @Accept       mpfd
// @Produce      json
// @Security     BearerAuth
// @Param        period_public_id path string true "Period public ID"
// @Param        file formData file true "Excel file (.xlsx)"
// @Success      200 {object} response.Response{data=model.ImportParticipantResponse}
// @Failure      400 {object} response.ResponseNoData
// @Failure      404 {object} response.ResponseNoData
// @Router       /periods/{period_public_id}/participants/import [post]
func (ctrl *ParticipantController) Import(c fiber.Ctx) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "excel file is required (field name: file)")
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if contentType != "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
		return fiber.NewError(fiber.StatusBadRequest, "file must be an Excel (.xlsx) file")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to read uploaded file")
	}
	defer file.Close()

	result, err := ctrl.UseCase.ImportFromExcel(c.Params("period_public_id"), file)
	if err != nil {
		return err
	}
	return response.Success(c, "Participants imported successfully.", result)
}