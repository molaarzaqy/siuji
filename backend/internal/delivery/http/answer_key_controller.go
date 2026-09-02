package http

import (
	"siuji-backend/internal/model"
	"siuji-backend/internal/usecase"
	"siuji-backend/pkg/response"

	"github.com/gofiber/fiber/v3"
)

type AnswerKeyController struct {
	UseCase *usecase.AnswerKeyUseCase
}


func NewAnswerKeyController(useCase *usecase.AnswerKeyUseCase) *AnswerKeyController {
	return &AnswerKeyController{
		UseCase: useCase,
	}
}

func (ctrl *AnswerKeyController) Upsert(c fiber.Ctx) error {
	request := new(model.UpsertAnswerKeyRequest)
	if err := c.Bind().Body(request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	result, err := ctrl.UseCase.Upsert(c.Params("question_public_id"), request)
	if err != nil {
		return err
	}
	return response.Success(c, "Answer key saved successfully.", result)
}