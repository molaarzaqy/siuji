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

var audioTypes = []string{"audio/mpeg", "audio/mp3", "audio/wav", "audio/x-wav"}
var imageTypes = []string{"image/jpeg", "image/png"}

func (ctrl *QuestionController) Create(c fiber.Ctx) error {
	request := &model.QuestionRequest{
		Question: c.FormValue("question"),
	}
	if passage := c.FormValue("passage"); passage != "" {
		request.Passage = &passage
	}

	audioFile, closeAudio, err := extractOptionalFile(c, "audio", audioTypes)
	if err != nil {
		return err
	}
	defer closeAudio()

	imageFile, closeImage, err := extractOptionalFile(c, "image", imageTypes)
	if err != nil {
		return err
	}
	defer closeImage()

	result, err := ctrl.UseCase.Create(c, c.Params("section_public_id"), request, audioFile, imageFile)
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
	request := &model.QuestionRequest{
		Question: c.FormValue("question"),
	}
	if passage := c.FormValue("passage"); passage != "" {
		request.Passage = &passage
	}

	audioFile, closeAudio, err := extractOptionalFile(c, "audio", audioTypes)
	if err != nil {
		return err
	}
	defer closeAudio()

	imageFile, closeImage, err := extractOptionalFile(c, "image", imageTypes)
	if err != nil {
		return err
	}
	defer closeImage()

	result, err := ctrl.UseCase.Update(c, c.Params("question_public_id"), request, audioFile, imageFile)
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