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

// Upsert godoc
// @Summary      Set correct answer for a question
// @Description  Create or update the answer key. The option must belong to the given question.
// @Tags         AnswerKey
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        question_public_id path string true "Question public ID"
// @Param        request body model.UpsertAnswerKeyRequest true "Correct option"
// @Success      200 {object} response.Response{data=model.AnswerKeyResponse}
// @Failure      400 {object} response.ResponseNoData
// @Failure      404 {object} response.ResponseNoData
// @Router       /questions/{question_public_id}/answer-key [put]
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