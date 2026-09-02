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

func (ctrl *OptionController) Delete(c fiber.Ctx) error {
	if err := ctrl.UseCase.Delete(c.Params("option_public_id")); err != nil {
		return err
	}
	return response.SuccessNoData(c, "Option deleted successfully.")
}

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