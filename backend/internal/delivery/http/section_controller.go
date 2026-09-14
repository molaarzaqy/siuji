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

// Create godoc
// @Summary      Create section
// @Description  Create a new section (question bank). section_type must be listening, structure, or reading.
// @Tags         Section
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body model.SectionRequest true "Section data"
// @Success      201 {object} response.Response{data=model.SectionResponse}
// @Failure      400 {object} response.ResponseNoData
// @Router       /sections [post]
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

// GetAll godoc
// @Summary      List sections
// @Tags         Section
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} response.Response{data=[]model.SectionResponse}
// @Router       /sections [get]
func (ctrl *SectionController) GetAll(c fiber.Ctx) error {
	result, err := ctrl.UseCase.GetAll()
	if err != nil {
		return err
	}
	return response.Success(c, "List sections retrieved successfully.", result)
}

// GetDetail godoc
// @Summary      Get section detail
// @Description  Get a section with its questions, options, and correct answers.
// @Tags         Section
// @Produce      json
// @Security     BearerAuth
// @Param        section_public_id path string true "Section public ID"
// @Success      200 {object} response.Response{data=model.SectionDetailResponse}
// @Failure      404 {object} response.ResponseNoData
// @Router       /sections/{section_public_id} [get]
func (ctrl *SectionController) GetDetail(c fiber.Ctx) error {
	result, err := ctrl.UseCase.GetDetail(c.Params("section_public_id"))
	if err != nil {
		return err
	}
	return response.Success(c, "Section detail retrieved successfully.", result)
}

// Update godoc
// @Summary      Update section
// @Tags         Section
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        section_public_id path string true "Section public ID"
// @Param        request body model.SectionRequest true "Section data"
// @Success      200 {object} response.Response{data=model.SectionResponse}
// @Failure      404 {object} response.ResponseNoData
// @Router       /sections/{section_public_id} [put]
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

// Delete godoc
// @Summary      Delete section
// @Tags         Section
// @Produce      json
// @Security     BearerAuth
// @Param        section_public_id path string true "Section public ID"
// @Success      200 {object} response.ResponseNoData
// @Failure      404 {object} response.ResponseNoData
// @Router       /sections/{section_public_id} [delete]
func (ctrl *SectionController) Delete(c fiber.Ctx) error {
	if err := ctrl.UseCase.Delete(c.Params("section_public_id")); err != nil {
		return err
	}
	return response.SuccessNoData(c, "Section deleted successfully.")
}