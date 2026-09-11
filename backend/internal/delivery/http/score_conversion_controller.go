package http

import (
	"strconv"

	"siuji-backend/internal/model"
	"siuji-backend/internal/usecase"
	"siuji-backend/pkg/response"

	"github.com/gofiber/fiber/v3"
)

type ScoreConversionController struct {
	UseCase *usecase.ScoreConversionUseCase
}

func NewScoreConversionController(useCase *usecase.ScoreConversionUseCase) *ScoreConversionController {
	return &ScoreConversionController{UseCase: useCase}
}

func (ctrl *ScoreConversionController) Create(c fiber.Ctx) error {
	request := new(model.ScoreConversionRequest)
	if err := c.Bind().Body(request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	result, err := ctrl.UseCase.Create(request)
	if err != nil {
		return err
	}
	return response.Created(c, "Score conversion created successfully.", result)
}

func (ctrl *ScoreConversionController) GetAll(c fiber.Ctx) error {
	sectionType := c.Query("section_type", "")
	result, err := ctrl.UseCase.GetAll(sectionType)
	if err != nil {
		return err
	}
	return response.Success(c, "List score conversions retrieved successfully.", result)
}

func (ctrl *ScoreConversionController) Update(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}

	request := new(model.ScoreConversionRequest)
	if err := c.Bind().Body(request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	result, err := ctrl.UseCase.Update(uint(id), request)
	if err != nil {
		return err
	}
	return response.Success(c, "Score conversion updated successfully.", result)
}

func (ctrl *ScoreConversionController) Delete(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}

	if err := ctrl.UseCase.Delete(uint(id)); err != nil {
		return err
	}
	return response.SuccessNoData(c, "Score conversion deleted successfully.")
}

func (ctrl *ScoreConversionController) BulkCreate(c fiber.Ctx) error {
	request := new(model.BulkScoreConversionRequest)
	if err := c.Bind().Body(request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	result, err := ctrl.UseCase.BulkCreate(request)
	if err != nil {
		return err
	}
	return response.Created(c, "Score conversions saved successfully.", result)
}