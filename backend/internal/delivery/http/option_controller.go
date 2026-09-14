package http

import (
	"siuji-backend/internal/model"
	"siuji-backend/internal/usecase"
	"siuji-backend/pkg/response"

	"github.com/gofiber/fiber/v3"
)

type OptionController struct {
	UseCase *usecase.OptionUseCase
}

func NewOptionController(useCase *usecase.OptionUseCase) *OptionController {
	return &OptionController{
		UseCase: useCase,
	}
}

// Create godoc
// @Summary      Create option
// @Description  Create a new answer option for a question. Label (A, B, C, ...) is auto-derived from position.
// @Tags         Option
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        question_public_id path string true "Question public ID"
// @Param        request body model.OptionRequest true "Option data"
// @Success      201 {object} response.Response{data=model.OptionResponse}
// @Failure      404 {object} response.ResponseNoData
// @Router       /questions/{question_public_id}/options [post]
func (ctrl *OptionController) Create(c fiber.Ctx) error {
	request := new(model.OptionRequest)
	if err := c.Bind().Body(request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	result, err := ctrl.UseCase.Create(c.Params("question_public_id"), request)
	if err != nil {
		return err
	}
	return response.Created(c, "Option created successfully.", result)
}

// Update godoc
// @Summary      Update option
// @Tags         Option
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        option_public_id path string true "Option public ID"
// @Param        request body model.OptionRequest true "Option data"
// @Success      200 {object} response.Response{data=model.OptionResponse}
// @Failure      404 {object} response.ResponseNoData
// @Router       /options/{option_public_id} [put]
func (ctrl *OptionController) Update(c fiber.Ctx) error {
	request := new(model.OptionRequest)
	if err := c.Bind().Body(request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	result, err := ctrl.UseCase.Update(c.Params("option_public_id"), request)
	if err != nil {
		return err
	}
	return response.Success(c, "Option updated successfully.", result)
}

// Delete godoc
// @Summary      Delete option
// @Tags         Option
// @Produce      json
// @Security     BearerAuth
// @Param        option_public_id path string true "Option public ID"
// @Success      200 {object} response.ResponseNoData
// @Failure      404 {object} response.ResponseNoData
// @Router       /options/{option_public_id} [delete]
func (ctrl *OptionController) Delete(c fiber.Ctx) error {
	if err := ctrl.UseCase.Delete(c.Params("option_public_id")); err != nil {
		return err
	}
	return response.SuccessNoData(c, "Option deleted successfully.")
}

// Reorder godoc
// @Summary      Reorder options within a question
// @Description  Labels (A, B, C, ...) are automatically reassigned to match the new order.
// @Tags         Option
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        question_public_id path string true "Question public ID"
// @Param        request body model.ReorderOptionsRequest true "Ordered list of option public IDs"
// @Success      200 {object} response.ResponseNoData
// @Failure      400 {object} response.ResponseNoData
// @Router       /questions/{question_public_id}/options/reorder [put]
func (ctrl *OptionController) Reorder(c fiber.Ctx) error {
	request := new(model.ReorderOptionsRequest)
	if err := c.Bind().Body(request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if err := ctrl.UseCase.Reorder(request); err != nil {
		return err
	}
	return response.SuccessNoData(c, "Option positions updated successfully.")
}