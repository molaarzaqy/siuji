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


// Create godoc
// @Summary      Create question
// @Description  Create a new question in a section. audio and image files are optional depending on section type.
// @Tags         Question
// @Accept       mpfd
// @Produce      json
// @Security     BearerAuth
// @Param        section_public_id path string true "Section public ID"
// @Param        question formData string true "Question text"
// @Param        passage formData string false "Reading passage (for Reading section)"
// @Param        audio formData file false "Audio file (for Listening section)"
// @Param        image formData file false "Image file (optional)"
// @Success      201 {object} response.Response{data=model.QuestionResponse}
// @Failure      400 {object} response.ResponseNoData
// @Failure      404 {object} response.ResponseNoData
// @Router       /sections/{section_public_id}/questions [post]
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

// GetDetail godoc
// @Summary      Get question detail
// @Description  Get a question with its options and correct answer.
// @Tags         Question
// @Produce      json
// @Security     BearerAuth
// @Param        question_public_id path string true "Question public ID"
// @Success      200 {object} response.Response{data=model.QuestionDetailResponse}
// @Failure      404 {object} response.ResponseNoData
// @Router       /questions/{question_public_id} [get]
func (ctrl *QuestionController) GetDetail(c fiber.Ctx) error {
	result, err := ctrl.UseCase.GetDetail(c.Params("question_public_id"))
	if err != nil {
		return err
	}
	return response.Success(c, "Question detail retrieved successfully.", result)
}

// Update godoc
// @Summary      Update question
// @Description  Update question text/passage. audio/image are optional — omit to keep the existing files.
// @Tags         Question
// @Accept       mpfd
// @Produce      json
// @Security     BearerAuth
// @Param        question_public_id path string true "Question public ID"
// @Param        question formData string true "Question text"
// @Param        passage formData string false "Reading passage"
// @Param        audio formData file false "New audio file (optional)"
// @Param        image formData file false "New image file (optional)"
// @Success      200 {object} response.Response{data=model.QuestionResponse}
// @Failure      404 {object} response.ResponseNoData
// @Router       /questions/{question_public_id} [put]
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

// Delete godoc
// @Summary      Delete question
// @Tags         Question
// @Produce      json
// @Security     BearerAuth
// @Param        question_public_id path string true "Question public ID"
// @Success      200 {object} response.ResponseNoData
// @Failure      404 {object} response.ResponseNoData
// @Router       /questions/{question_public_id} [delete]
func (ctrl *QuestionController) Delete(c fiber.Ctx) error {
	if err := ctrl.UseCase.Delete(c.Params("question_public_id")); err != nil {
		return err
	}
	return response.SuccessNoData(c, "Question deleted successfully.")
}

// Reorder godoc
// @Summary      Reorder questions within a section
// @Tags         Question
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        section_public_id path string true "Section public ID"
// @Param        request body model.ReorderQuestionsRequest true "Ordered list of question public IDs"
// @Success      200 {object} response.ResponseNoData
// @Failure      400 {object} response.ResponseNoData
// @Router       /sections/{section_public_id}/questions/reorder [put]
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