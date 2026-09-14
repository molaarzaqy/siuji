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

// Create godoc
// @Summary      Create score conversion entry
// @Tags         ScoreConversion
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body model.ScoreConversionRequest true "Conversion entry"
// @Success      201 {object} response.Response{data=model.ScoreConversionResponse}
// @Failure      400 {object} response.ResponseNoData
// @Router       /score-conversions [post]
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

// GetAll godoc
// @Summary      List score conversions
// @Tags         ScoreConversion
// @Produce      json
// @Security     BearerAuth
// @Param        section_type query string false "Filter by listening, structure, or reading"
// @Success      200 {object} response.Response{data=[]model.ScoreConversionResponse}
// @Router       /score-conversions [get]
func (ctrl *ScoreConversionController) GetAll(c fiber.Ctx) error {
	sectionType := c.Query("section_type", "")
	result, err := ctrl.UseCase.GetAll(sectionType)
	if err != nil {
		return err
	}
	return response.Success(c, "List score conversions retrieved successfully.", result)
}

// Update godoc
// @Summary      Update score conversion entry
// @Tags         ScoreConversion
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Score conversion ID"
// @Param        request body model.ScoreConversionRequest true "Conversion entry"
// @Success      200 {object} response.Response{data=model.ScoreConversionResponse}
// @Failure      404 {object} response.ResponseNoData
// @Router       /score-conversions/{id} [put]
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

// Delete godoc
// @Summary      Delete score conversion entry
// @Tags         ScoreConversion
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Score conversion ID"
// @Success      200 {object} response.ResponseNoData
// @Failure      404 {object} response.ResponseNoData
// @Router       /score-conversions/{id} [delete]
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

// BulkCreate godoc
// @Summary      Bulk upsert score conversions
// @Description  Insert or update multiple entries at once, keyed by (section_type, correct_count).
// @Tags         ScoreConversion
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body model.BulkScoreConversionRequest true "List of conversion entries"
// @Success      201 {object} response.Response{data=model.BulkScoreConversionResponse}
// @Failure      400 {object} response.ResponseNoData
// @Router       /score-conversions/bulk [post]
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