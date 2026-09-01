package http

import (
	"strconv"

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

func (ctrl *PeriodController) Create(c fiber.Ctx) error {
	request := new(model.PeriodRequest)
	if err := c.Bind().Body(request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	result, err := ctrl.UseCase.Create(request)
	if err != nil {
		return err
	}
	return response.Created(c, "Period created successfully.", result)
}

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

func (ctrl *PeriodController) GetDetail(c fiber.Ctx) error {
	result, err := ctrl.UseCase.GetDetail(c.Params("period_public_id"))
	if err != nil {
		return err
	}
	return response.Success(c, "Period detail retrieved successfully.", result)
}

func (ctrl *PeriodController) Update(c fiber.Ctx) error {
	request := new(model.PeriodRequest)
	if err := c.Bind().Body(request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	result, err := ctrl.UseCase.Update(c.Params("period_public_id"), request)
	if err != nil {
		return err
	}
	return response.Success(c, "Period updated successfully.", result)
}

func (ctrl *PeriodController) Delete(c fiber.Ctx) error {
	if err := ctrl.UseCase.Delete(c.Params("period_public_id")); err != nil {
		return err
	}
	return response.SuccessNoData(c, "Period deleted successfully.")
}

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

func (ctrl *PeriodController) RemoveSection(c fiber.Ctx) error {
	err := ctrl.UseCase.RemoveSection(c.Params("period_public_id"), c.Params("section_public_id"))
	if err != nil {
		return err
	}
	return response.SuccessNoData(c, "Section removed from period successfully.")
}

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