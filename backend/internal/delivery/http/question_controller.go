package http

import (
	"siuji-backend/internal/model"
	"siuji-backend/internal/usecase"
	"siuji-backend/pkg/response"

	"github.com/gofiber/fiber/v3"
)

type QuestionController struct {
	UseCase *usecase.QuestionUseCase
}

func NewQuestionController(useCase *usecase.QuestionUseCase) *QuestionController {
	return &QuestionController{
		UseCase: useCase,
	}
}

func (ctrl *QuestionController) Create(c fiber.Ctx) error {
	request := new(model.QuestionRequest)
	if err := c.Bind().Body(request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	result, err := ctrl.UseCase.Create(c.Params("section_public_id"), request)
	if err != nil {
		return err
	}
	return response.Created(c, "Question created successfully.", result)
}

func (ctrl *QuestionController) GetDetail(c fiber.Ctx) error {
	result, err := ctrl.UseCase.GetDetail(c.Params("question_public_id"))
	if err != nil {
		return err
	}
	return response.Success(c, "Question detail retrieved successfully.", result)
}

func (ctrl *QuestionController) Update(c fiber.Ctx) error {
	request := new(model.QuestionRequest)
	if err := c.Bind().Body(request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	result, err := ctrl.UseCase.Update(c.Params("question_public_id"), request)
	if err != nil {
		return err
	}
	return response.Success(c, "Question updated successfully.", result)
}

func (ctrl *QuestionController) Delete(c fiber.Ctx) error {
	if err := ctrl.UseCase.Delete(c.Params("question_public_id")); err != nil {
		return err
	}
	return response.SuccessNoData(c, "Question deleted successfully.")
}

func (ctrl *QuestionController) Reorder(c fiber.Ctx) error {
	request := new(model.ReorderQuestionsRequest)
	if err := c.Bind().Body(request); err != nil {
		return err
	}
	if err := ctrl.UseCase.Reorder(request); err != nil {
		return err
	}
	return response.SuccessNoData(c, "Question positions updated successfully.")
}