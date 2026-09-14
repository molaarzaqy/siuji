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

// GetPeriods godoc
// @Summary      List my exam periods
// @Description  Get periods the logged-in participant is registered for.
// @Tags         ParticipantExam
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} response.Response{data=[]model.ParticipantPeriodListResponse}
// @Router       /participant/periods [get]
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

// GetPeriodDetail godoc
// @Summary      Get my exam period detail
// @Tags         ParticipantExam
// @Produce      json
// @Security     BearerAuth
// @Param        period_public_id path string true "Period public ID"
// @Success      200 {object} response.Response{data=model.ParticipantPeriodDetailResponse}
// @Failure      404 {object} response.ResponseNoData
// @Router       /participant/periods/{period_public_id} [get]
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

// StartExam godoc
// @Summary      Start exam
// @Description  Marks the exam as started and returns all sections/questions (without correct answers).
// @Tags         ParticipantExam
// @Produce      json
// @Security     BearerAuth
// @Param        period_public_id path string true "Period public ID"
// @Success      200 {object} response.Response{data=model.ExamSessionResponse}
// @Failure      400 {object} response.ResponseNoData
// @Failure      404 {object} response.ResponseNoData
// @Router       /participant/periods/{period_public_id}/start [post]
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

// SaveAnswer godoc
// @Summary      Save/update an answer
// @Description  Can be called repeatedly while the exam is in progress — upserts by question.
// @Tags         ParticipantExam
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        period_public_id path string true "Period public ID"
// @Param        request body model.SaveAnswerRequest true "Answer"
// @Success      200 {object} response.Response{data=model.SaveAnswerResponse}
// @Failure      400 {object} response.ResponseNoData
// @Failure      404 {object} response.ResponseNoData
// @Router       /participant/periods/{period_public_id}/answers [post]
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

// SubmitExam godoc
// @Summary      Submit exam
// @Description  Finalizes the exam, computes section scores and the final TOEFL score.
// @Tags         ParticipantExam
// @Produce      json
// @Security     BearerAuth
// @Param        period_public_id path string true "Period public ID"
// @Success      200 {object} response.Response{data=model.SubmitExamResponse}
// @Failure      400 {object} response.ResponseNoData
// @Failure      404 {object} response.ResponseNoData
// @Router       /participant/periods/{period_public_id}/submit [post]
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

// GetResult godoc
// @Summary      Get exam result
// @Description  Returns the final score, section breakdown, and certificate URL (generated on first request if passed).
// @Tags         ParticipantExam
// @Produce      json
// @Security     BearerAuth
// @Param        period_public_id path string true "Period public ID"
// @Success      200 {object} response.Response{data=model.ExamResultResponse}
// @Failure      404 {object} response.ResponseNoData
// @Router       /participant/periods/{period_public_id}/result [get]
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

