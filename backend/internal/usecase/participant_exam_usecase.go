package usecase

import (
	"bytes"
	"context"
	"time"

	"siuji-backend/internal/entity"
	"siuji-backend/internal/model"
	"siuji-backend/internal/model/converter"
	"siuji-backend/internal/repository"
	"siuji-backend/pkg/certificate"
	"siuji-backend/pkg/cloudinary"
	"siuji-backend/pkg/scorer"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
)

type ParticipantExamUseCase struct {
	Log                         *logrus.Logger
	Validate                    *validator.Validate
	ParticipantPeriodRepository repository.ParticipantPeriodRepository
	PeriodRepository            repository.PeriodRepository
	PeriodSectionRepository     repository.PeriodSectionRepository
	SectionRepository           repository.SectionRepository
	QuestionRepository          repository.QuestionRepository
	OptionRepository            repository.OptionRepository
	AnswerKeyRepository         repository.AnswerKeyRepository
	ParticipantAnswerRepository repository.ParticipantAnswerRepository
	SectionScoreRepository      repository.SectionScoreRepository
	ScoreConversionRepository	repository.ScoreConversionRepository
	CloudinaryService           *cloudinary.Service
	CertificateGenerator        *certificate.Generator
}

func NewParticipantExamUseCase(
	log *logrus.Logger,
	validate *validator.Validate,
	participantPeriodRepo repository.ParticipantPeriodRepository,
	periodRepo repository.PeriodRepository,
	periodSectionRepo repository.PeriodSectionRepository,
	sectionRepo repository.SectionRepository,
	questionRepo repository.QuestionRepository,
	optionRepo repository.OptionRepository,
	answerKeyRepo repository.AnswerKeyRepository,
	participantAnswerRepo repository.ParticipantAnswerRepository,
	sectionScoreRepo repository.SectionScoreRepository,
	scoreConversionRepository	repository.ScoreConversionRepository,
	cloudinaryService *cloudinary.Service,
	certGenerator *certificate.Generator,
) *ParticipantExamUseCase {
	return &ParticipantExamUseCase{
		Log:                         log,
		Validate:                    validate,
		ParticipantPeriodRepository: participantPeriodRepo,
		PeriodRepository:            periodRepo,
		PeriodSectionRepository:     periodSectionRepo,
		SectionRepository:           sectionRepo,
		QuestionRepository:          questionRepo,
		OptionRepository:            optionRepo,
		AnswerKeyRepository:         answerKeyRepo,
		ParticipantAnswerRepository: participantAnswerRepo,
		SectionScoreRepository:      sectionScoreRepo,
		ScoreConversionRepository: scoreConversionRepository,
		CloudinaryService:           cloudinaryService,
		CertificateGenerator:        certGenerator,
	}
}

// 1. GET /api/v1/participant/periods
func (u *ParticipantExamUseCase) GetPeriods(userID uint) ([]model.ParticipantPeriodListResponse, error) {
	pps, err := u.ParticipantPeriodRepository.FindByUserID(userID)
	if err != nil {
		u.Log.Errorf("failed to fetch participant periods for user %d: %v", userID, err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to fetch participant periods")
	}

	responses := make([]model.ParticipantPeriodListResponse, 0, len(pps))
	for _, pp := range pps {
		responses = append(responses, converter.ParticipantPeriodToListResponse(&pp))
	}
	return responses, nil
}

// 2. GET /api/v1/participant/periods/:period_public_id
func (u *ParticipantExamUseCase) GetPeriodDetail(userID uint, periodPublicID string) (*model.ParticipantPeriodDetailResponse, error) {
	pp, err := u.ParticipantPeriodRepository.FindByPeriodPublicIDAndUserID(periodPublicID, userID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "you are not registered for this period or period not found")
	}

	return converter.PeriodToParticipantDetailResponse(&pp.Period, pp.Status), nil
}

// 3. POST /api/v1/participant/periods/:period_public_id/start
func (u *ParticipantExamUseCase) StartExam(userID uint, periodPublicID string) (*model.ExamSessionResponse, error) {
	pp, err := u.ParticipantPeriodRepository.FindByPeriodPublicIDAndUserID(periodPublicID, userID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "you are not registered for this period or period not found")
	}

	if pp.Status == "completed" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "you have already completed this exam")
	}

	now := time.Now()
	if !pp.Period.StartTime.IsZero() && now.Before(pp.Period.StartTime) {
		return nil, fiber.NewError(fiber.StatusBadRequest, "exam period has not started yet")
	}
	if !pp.Period.EndTime.IsZero() && now.After(pp.Period.EndTime) {
		return nil, fiber.NewError(fiber.StatusBadRequest, "exam period has already ended")
	}

	// Update status from registered to started
	if pp.Status == "registered" {
		pp.Status = "started"
		pp.UpdatedAt = now
		if err := u.ParticipantPeriodRepository.Update(pp); err != nil {
			u.Log.Errorf("failed to update participant status to started: %v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to start exam")
		}
	}

	// Fetch sections in this period
	periodSections, err := u.PeriodSectionRepository.FindByPeriodID(pp.PeriodID)
	if err != nil {
		u.Log.Errorf("failed to fetch period sections: %v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to load exam sections")
	}

	examSections := make([]model.ExamSectionResponse, 0, len(periodSections))
	for _, ps := range periodSections {
		// Fetch section with its questions and options
		sectionWithQuestions, err := u.SectionRepository.FindByPublicIDWithQuestions(ps.Section.PublicID.String())
		if err != nil {
			u.Log.Warnf("failed to load questions for section %s: %v", ps.Section.PublicID, err)
			continue
		}

		examQuestions := make([]model.ExamQuestionResponse, 0, len(sectionWithQuestions.Questions))
		for _, q := range sectionWithQuestions.Questions {
			examQuestions = append(examQuestions, converter.QuestionToExamQuestionResponse(&q))
		}

		examSections = append(examSections, model.ExamSectionResponse{
			PublicID:  ps.Section.PublicID.String(),
			Title:     ps.Section.Title,
			Position:  ps.Position,
			Questions: examQuestions,
		})
	}

	return &model.ExamSessionResponse{
		PeriodPublicID: pp.Period.PublicID.String(),
		Title:          pp.Period.Title,
		Sections:       examSections,
	}, nil
}

// 4. POST /api/v1/participant/periods/:period_public_id/answers
func (u *ParticipantExamUseCase) SaveAnswer(userID uint, periodPublicID string, req *model.SaveAnswerRequest) (*model.SaveAnswerResponse, error) {
	if err := u.Validate.Struct(req); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	pp, err := u.ParticipantPeriodRepository.FindByPeriodPublicIDAndUserID(periodPublicID, userID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "you are not registered for this period or period not found")
	}

	if pp.Status != "started" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "exam is not currently in progress")
	}

	now := time.Now()
	if !pp.Period.EndTime.IsZero() && now.After(pp.Period.EndTime) {
		return nil, fiber.NewError(fiber.StatusBadRequest, "exam period has already ended")
	}

	// Verify question exists and belongs to a section of this period
	question, err := u.QuestionRepository.FindByPublicID(req.QuestionPublicID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "question not found")
	}

	isSectionInPeriod, err := u.PeriodSectionRepository.ExistsByPeriodAndSection(pp.PeriodID, question.SectionID)
	if err != nil || !isSectionInPeriod {
		return nil, fiber.NewError(fiber.StatusBadRequest, "question does not belong to this exam period")
	}

	// Verify option exists and belongs to the question
	option, err := u.OptionRepository.FindByPublicID(req.OptionPublicID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "option not found")
	}
	if option.QuestionID != question.ID {
		return nil, fiber.NewError(fiber.StatusBadRequest, "option does not belong to the specified question")
	}

	// Check if this option is correct
	answerKey, err := u.AnswerKeyRepository.FindByQuestionID(question.ID)
	if err != nil {
		u.Log.Warnf("failed to lookup answer key for question %d: %v", question.ID, err)
	}
	isCorrect := (answerKey != nil && answerKey.OptionID == option.ID)

	answer := &entity.ParticipantAnswer{
		ParticipantPeriodID: pp.ID,
		QuestionID:          question.ID,
		OptionID:            option.ID,
		IsCorrect:           isCorrect,
	}

	if err := u.ParticipantAnswerRepository.Upsert(answer); err != nil {
		u.Log.Errorf("failed to upsert participant answer: %v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to save answer")
	}

	return &model.SaveAnswerResponse{
		QuestionPublicID: req.QuestionPublicID,
		OptionPublicID:   req.OptionPublicID,
		UpdatedAt:        answer.UpdatedAt,
	}, nil
}

// 5. POST /api/v1/participant/periods/:period_public_id/submit
func (u *ParticipantExamUseCase) SubmitExam(userID uint, periodPublicID string) (*model.SubmitExamResponse, error) {
	pp, err := u.ParticipantPeriodRepository.FindByPeriodPublicIDAndUserID(periodPublicID, userID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "you are not registered for this period or period not found")
	}
	if pp.Status == "completed" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "exam has already been submitted")
	}
	if pp.Status != "started" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "exam has not been started yet")
	}
	periodSections, err := u.PeriodSectionRepository.FindByPeriodID(pp.PeriodID)
	if err != nil {
		u.Log.Errorf("failed to fetch period sections for scoring: %v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to grade exam")
	}
	listeningScaled, structureScaled, readingScaled := 0, 0, 0

		for _, ps := range periodSections {
		correctCount, err := u.ParticipantAnswerRepository.CountCorrectByParticipantPeriodAndSection(pp.ID, ps.SectionID)
		if err != nil {
			u.Log.Errorf("failed to count correct answers for section %d: %v", ps.SectionID, err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to grade exam")
		}

		sectionType := ps.Section.SectionType
		if sectionType == "" {
			u.Log.Errorf("section %d has no section_type configured, cannot grade accurately", ps.SectionID)
			return nil, fiber.NewError(fiber.StatusInternalServerError, "exam configuration incomplete: section type not set for one or more sections")
		}

		scaledScore, err := u.ScoreConversionRepository.FindScaledScore(sectionType, correctCount)
		if err != nil {
			u.Log.Errorf("failed to find score conversion for %s (%d correct): %v", sectionType, correctCount, err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to grade exam: score conversion table incomplete for "+sectionType)
		}

		switch sectionType {
		case "listening":
			listeningScaled = scaledScore
		case "structure":
			structureScaled = scaledScore
		case "reading":
			readingScaled = scaledScore
		}
				secScore := &entity.SectionScore{
			ParticipantPeriodID: pp.ID,
			SectionID:           ps.SectionID,
			CorrectCount:        correctCount,
			RawScore:            correctCount,
			ScaledScore:         scaledScore,
		}
		if err := u.SectionScoreRepository.Upsert(secScore); err != nil {
			u.Log.Errorf("failed to save section score for section %d: %v", ps.SectionID, err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to save exam scores")
		}
	}

	finalScore := scorer.CalculateFinalTOEFL(listeningScaled, structureScaled, readingScaled)

	now := time.Now()
	pp.Status = "completed"
	pp.Score = &finalScore
	pp.UpdatedAt = now

	if err := u.ParticipantPeriodRepository.Update(pp); err != nil {
		u.Log.Errorf("failed to update participant period status to completed: %v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to submit exam")
	}

	return &model.SubmitExamResponse{
		PeriodPublicID: periodPublicID,
		Status:         "completed",
		SubmittedAt:    now,
	}, nil

}

// 6. GET /api/v1/participant/periods/:period_public_id/result
func (u *ParticipantExamUseCase) GetResult(ctx context.Context, userID uint, periodPublicID string) (*model.ExamResultResponse, error) {
	pp, err := u.ParticipantPeriodRepository.FindByPeriodPublicIDAndUserID(periodPublicID, userID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "you are not registered for this period or period not found")
	}

	if pp.Status != "completed" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "exam has not been submitted yet")
	}

	finalScore := 0
	if pp.Score != nil {
		finalScore = *pp.Score
	}

	status := "failed"
	if finalScore >= pp.Period.MinPassingGrade {
		status = "passed"
	}

	// Fetch section breakdowns
	sectionScores, err := u.SectionScoreRepository.FindByParticipantPeriodID(pp.ID)
	if err != nil {
		u.Log.Errorf("failed to fetch section scores: %v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to retrieve score breakdowns")
	}

	breakdowns := make([]model.SectionBreakdownResponse, 0, len(sectionScores))
	listeningScore := 0
	structureScore := 0
	readingScore := 0

	for _, ss := range sectionScores {
		breakdowns = append(breakdowns, model.SectionBreakdownResponse{
			SectionTitle: ss.Section.Title,
			CorrectCount: ss.CorrectCount,
			ScaledScore:  ss.ScaledScore,
		})
	
		switch ss.Section.SectionType {
		case "listening":
			listeningScore = ss.ScaledScore
		case "structure":
			structureScore = ss.ScaledScore
		case "reading":
			readingScore = ss.ScaledScore
		}
	}

	var certificateURL *string
	if pp.CertificateURL != nil && *pp.CertificateURL != "" {
		certificateURL = pp.CertificateURL
	} else if status == "passed" {
		// Generate Certificate PDF
		expFormatted := ""
		if !pp.Period.CertificateExpMonth.IsZero() {
			expFormatted = pp.Period.CertificateExpMonth.Format("02 January 2006")
		}

		certData := &certificate.CertificateData{
			ParticipantName: pp.User.Name,
			ParticipantNIM:  pp.User.NIM,
			University:      pp.User.University,
			PeriodTitle:     pp.Period.Title,
			CertificateExp:  expFormatted,
			ListeningScore:  listeningScore,
			StructureScore:  structureScore,
			ReadingScore:    readingScore,
			FinalScore:      finalScore,
			TemplateURL:     pp.Period.CertificateURL,
		}

		pdfBytes, err := u.CertificateGenerator.Generate(certData)
		if err != nil {
			u.Log.Errorf("failed to generate certificate pdf: %v", err)
		} else {
			uploadedURL, err := u.CloudinaryService.UploadGeneratedCertificate(ctx, bytes.NewReader(pdfBytes))
			if err != nil {
				u.Log.Errorf("failed to upload generated certificate to cloudinary: %v", err)
			} else {
				certificateURL = &uploadedURL
				pp.CertificateURL = &uploadedURL
				if err := u.ParticipantPeriodRepository.Update(pp); err != nil {
					u.Log.Errorf("failed to save certificate url to participant_period: %v", err)
				}
			}
		}
	}

	return &model.ExamResultResponse{
		PeriodTitle:       pp.Period.Title,
		Status:            status,
		FinalScore:        finalScore,
		MinPassingGrade:   pp.Period.MinPassingGrade,
		CertificateURL:    certificateURL,
		SectionBreakdowns: breakdowns,
	}, nil
}

