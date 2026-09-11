package http

import (
	"siuji-backend/internal/model"
	"siuji-backend/internal/usecase"
	"siuji-backend/pkg/response"

	"github.com/gofiber/fiber/v3"
)

type ParticipantExamController struct {
	UseCase *usecase.ParticipantExamUseCase
}

func NewParticipantExamController(useCase *usecase.ParticipantExamUseCase) *ParticipantExamController {
	return &ParticipantExamController{UseCase: useCase}
}

func (ctrl *ParticipantExamController) getUserID(c fiber.Ctx) (uint, error) {
	if val := c.Locals("user_id"); val != nil {
		switch v := val.(type) {
		case uint:
			return v, nil
		case int:
			return uint(v), nil
		case float64:
			return uint(v), nil
		}
	}
	return 0, fiber.NewError(fiber.StatusUnauthorized, "unauthorized: user id not found in token")
}

// 1. GET /api/v1/participant/periods
func (ctrl *ParticipantExamController) GetPeriods(c fiber.Ctx) error {
	userID, err := ctrl.getUserID(c)
	if err != nil {
		return err
	}

	result, err := ctrl.UseCase.GetPeriods(userID)
	if err != nil {
		return err
	}
	return response.Success(c, "Participant periods retrieved successfully.", result)
}

// 2. GET /api/v1/participant/periods/:period_public_id
func (ctrl *ParticipantExamController) GetPeriodDetail(c fiber.Ctx) error {
	userID, err := ctrl.getUserID(c)
	if err != nil {
		return err
	}

	result, err := ctrl.UseCase.GetPeriodDetail(userID, c.Params("period_public_id"))
	if err != nil {
		return err
	}
	return response.Success(c, "Period detail retrieved successfully.", result)
}

// 3. POST /api/v1/participant/periods/:period_public_id/start
func (ctrl *ParticipantExamController) StartExam(c fiber.Ctx) error {
	userID, err := ctrl.getUserID(c)
	if err != nil {
		return err
	}

	result, err := ctrl.UseCase.StartExam(userID, c.Params("period_public_id"))
	if err != nil {
		return err
	}
	return response.Success(c, "Exam session started successfully.", result)
}

// 4. POST /api/v1/participant/periods/:period_public_id/answers
func (ctrl *ParticipantExamController) SaveAnswer(c fiber.Ctx) error {
	userID, err := ctrl.getUserID(c)
	if err != nil {
		return err
	}

	req := new(model.SaveAnswerRequest)
	if err := c.Bind().Body(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	result, err := ctrl.UseCase.SaveAnswer(userID, c.Params("period_public_id"), req)
	if err != nil {
		return err
	}
	return response.Success(c, "Answer saved successfully.", result)
}

// 5. POST /api/v1/participant/periods/:period_public_id/submit
func (ctrl *ParticipantExamController) SubmitExam(c fiber.Ctx) error {
	userID, err := ctrl.getUserID(c)
	if err != nil {
		return err
	}

	result, err := ctrl.UseCase.SubmitExam(userID, c.Params("period_public_id"))
	if err != nil {
		return err
	}
	return response.Success(c, "Exam submitted successfully.", result)
}

// 6. GET /api/v1/participant/periods/:period_public_id/result
func (ctrl *ParticipantExamController) GetResult(c fiber.Ctx) error {
	userID, err := ctrl.getUserID(c)
	if err != nil {
		return err
	}

	result, err := ctrl.UseCase.GetResult(c.Context(), userID, c.Params("period_public_id"))
	if err != nil {
		return err
	}
	return response.Success(c, "Exam result retrieved successfully.", result)
}

