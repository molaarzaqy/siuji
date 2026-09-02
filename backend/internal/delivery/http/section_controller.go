package http

import (
	"siuji-backend/internal/model"
	"siuji-backend/internal/usecase"
	"siuji-backend/pkg/response"

	"github.com/gofiber/fiber/v3"
)

type SectionController struct {
	UseCase *usecase.SectionUseCase
}

func NewSectionController(useCase *usecase.SectionUseCase) *SectionController {
	return &SectionController{UseCase: useCase}
}

func (ctrl *SectionController) Create(c fiber.Ctx) error {
	request := new(model.SectionRequest)
	if err := c.Bind().Body(request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	result, err := ctrl.UseCase.Create(request)
	if err != nil {
		return err
	}
	return response.Created(c, "Section created successfully.", result)
}

func (ctrl *SectionController) GetAll(c fiber.Ctx) error {
	result, err := ctrl.UseCase.GetAll()
	if err != nil {
		return err
	}
	return response.Success(c, "List sections retrieved successfully.", result)
}

func (ctrl *SectionController) GetDetail(c fiber.Ctx) error {
	result, err := ctrl.UseCase.GetDetail(c.Params("section_public_id"))
	if err != nil {
		return err
	}
	return response.Success(c, "Section detail retrieved successfully.", result)
}

func (ctrl *SectionController) Update(c fiber.Ctx) error {
	request := new(model.SectionRequest)
	if err := c.Bind().Body(request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	result, err := ctrl.UseCase.Update(c.Params("section_public_id"), request)
	if err != nil {
		return err
	}
	return response.Success(c, "Section updated successfully.", result)
}

func (ctrl *SectionController) Delete(c fiber.Ctx) error {
	if err := ctrl.UseCase.Delete(c.Params("section_public_id")); err != nil {
		return err
	}
	return response.SuccessNoData(c, "Section deleted successfully.")
}