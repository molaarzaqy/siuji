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

func (ctrl *ParticipantController) GetDetail(c fiber.Ctx) error {
	result, err := ctrl.UseCase.GetDetail(c.Params("period_public_id"), c.Params("user_public_id"))
	if err != nil {
		return err
	}
	return response.Success(c, "Participant detail retrieved successfully.", result)
}

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

func (ctrl *ParticipantController) Remove(c fiber.Ctx) error {
	err := ctrl.UseCase.Remove(c.Params("period_public_id"), c.Params("user_public_id"))
	if err != nil {
		return err
	}
	return response.SuccessNoData(c, "Participant removed from period successfully.")
}