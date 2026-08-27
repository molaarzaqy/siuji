package services

import (
	"errors"
	"siuji-backend/models"
	"siuji-backend/repositories"
	"siuji-backend/utils"
	"time"
)

type ParticipantExamService interface {
	GetMyPeriods(userID uint) ([]models.ParticipantPeriodListResponse, error)
	GetMyPeriodDetail(periodPublicID string, userID uint) (*models.ParticipantPeriodDetailResponse, error)
	StartExam(periodPublicID string, userID uint) (*models.ExamStartResponse, error)
	SaveAnswer(periodPublicID string, userID uint, req *models.SaveAnswerRequest) (*models.AnswerResponse, error)
	SubmitExam(periodPublicID string, userID uint) (*models.ExamSubmitResponse, error)
}

type participantExamService struct {
	participantRepo repositories.ParticipantRepository
	periodRepo repositories.PeriodRepository
	answerKeyRepo repositories.AnswerKeyRepository
	participantAnswerRepo repositories.ParticipantAnswerRepository
	sectionScoreRepo repositories.SectionScoreRepository
	questionRepo repositories.QuestionRepository
	optionRepo repositories.OptionRepository
}

func NewPeriodExamService(
	participantRepo repositories.ParticipantRepository, 
	periodRepo repositories.PeriodRepository,
	answerKeyRepo repositories.AnswerKeyRepository,
	participantAnswerRepo repositories.ParticipantAnswerRepository,
	sectionScoreRepo repositories.SectionScoreRepository,
	questionRepo repositories.QuestionRepository,
	optionRepo repositories.OptionRepository,
) ParticipantExamService {
	return &participantExamService{
		participantRepo: participantRepo,
		periodRepo: periodRepo,
		answerKeyRepo: answerKeyRepo,
		participantAnswerRepo: participantAnswerRepo,
		sectionScoreRepo: sectionScoreRepo,
		questionRepo: questionRepo,
		optionRepo: optionRepo,
	}
}

func (s *participantExamService) GetMyPeriods(userID uint) ([]models.ParticipantPeriodListResponse, error) {
	participantPeriods, err := s.participantRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	var responses []models.ParticipantPeriodListResponse
	for _, pp := range participantPeriods {
		responses = append(responses, models.ParticipantPeriodListResponse{
			PublicID: pp.PublicID,
			Period: models.PeriodItemResponse{
				PublicID:  pp.Period.PublicID,
				Title:     pp.Period.Title,
				StartTime: pp.Period.StartTime,
				EndTime:   pp.Period.EndTime,
				Status:    pp.Period.Status,
			},
			Status: pp.Status,
			Score:  pp.Score,
		})
	}

	return responses, nil
}

func (s *participantExamService) GetMyPeriodDetail(periodPublicID string, userID uint) (*models.ParticipantPeriodDetailResponse, error) {
	// Validasi apakah periode dengan public_id tersebut ada
	period, err := s.periodRepo.FindByPublicID(periodPublicID)
	if err != nil || period == nil {
		return nil, errors.New("period not found")
	}
	// Cari relasi participant period berdasarkan internal period ID dan user ID
	participantPeriod, err := s.participantRepo.FindByPeriodIDAndUserID(period.ID, userID)
	if err != nil {
		return nil, err
	}
	if participantPeriod == nil {
		return nil, errors.New("you are not registered in this period")
	}
	// Mapping ke response DTO yang diharapkan
	response := &models.ParticipantPeriodDetailResponse{
		PublicID:          period.PublicID,
		Title:             period.Title,
		Month:             period.Month,
		Year:              period.Year,
		StartTime:         period.StartTime,
		EndTime:           period.EndTime,
		ParticipantStatus: participantPeriod.Status,
	}

	return response, nil
}

func (s *participantExamService) StartExam(periodPublicID string, userID uint) (*models.ExamStartResponse, error) {
	// Ambil data periode beserta soal lengkap lewat tabel pivot period_sections
	period, err := s.periodRepo.FindWithSectionsAndQuestions(periodPublicID)
	if err != nil || period == nil {
		return nil, errors.New("period not found")
	}
	// Validasi status periode (harus published)
	if period.Status != utils.PeriodStatusPublished {
		return nil, errors.New("exam period is not active")
	}
	// Validasi waktu ujian (Opsional tapi disarankan: cek apakah sekarang di antara StartTime dan EndTime)
	now := time.Now()
	if now.Before(period.StartTime) {
		return nil, errors.New("exam has not started yet")
	}
	if now.After(period.EndTime) {
		return nil, errors.New("exam session has already ended")
	}
	// Cari data participant period milik user yang bersangkutan
	participantPeriod, err := s.participantRepo.FindByPeriodIDAndUserID(period.ID, userID)
	if err != nil {
		return nil, err
	}
	if participantPeriod == nil {
		return nil, errors.New("you are not registered in this period")
	}
	// Cek status peserta, jika masih 'registered', ubah menjadi 'started'
	if participantPeriod.Status == utils.ParticipantStatusRegistered {
		err := s.participantRepo.UpdateStatus(participantPeriod.ID, utils.ParticipantStatusStarted)
		if err != nil {
			return nil, errors.New("failed to start exam session")
		}
	} else if participantPeriod.Status == utils.ParticipantStatusCompleted {
		return nil, errors.New("exam has already been completed")
	}
	// Mapping data database ke struktur DTO Response yang diinginkan
	var sectionResponses []models.SectionExamItem
	for _, ps := range period.PeriodSections {
		var questionResponses []models.QuestionExamItem
		for _, q := range ps.Section.Questions {
			var optionResponses []models.OptionExamItem
			for _, opt := range q.Options {
				optionResponses = append(optionResponses, models.OptionExamItem{
					PublicID:   opt.PublicID,
					Label:      opt.Label,
					OptionText: opt.OptionText,
					Position:   opt.Position,
				})
			}
			questionResponses = append(questionResponses, models.QuestionExamItem{
				PublicID:  q.PublicID,
				Question:  q.Question,
				AudioURL:  q.AudioURL,
				ImageURL:  q.ImageURL,
				Passage:   q.Passage,
				Position:  q.Position,
				Options:   optionResponses,
			})
		}
		sectionResponses = append(sectionResponses, models.SectionExamItem{
			PublicID:  ps.Section.PublicID,
			Title:     ps.Section.Title,
			Position:  ps.Position,
			Questions: questionResponses,
		})
	}
	response := &models.ExamStartResponse{
		PeriodPublicID: period.PublicID,
		Title:          period.Title,
		Sections:       sectionResponses,
	}

	return response, nil
}

func (s *participantExamService) SaveAnswer(periodPublicID string, userID uint, req *models.SaveAnswerRequest) (*models.AnswerResponse, error) {
	// Validasi periode dan waktu mutlak server
	period, err := s.periodRepo.FindByPublicID(periodPublicID)
	if err != nil || period == nil {
		return nil, errors.New("period not found")
	}

	now := time.Now()
	if now.After(period.EndTime) {
		return nil, errors.New("exam time is over, cannot save answer")
	}

	// Validasi status sesi participant
	participantPeriod, err := s.participantRepo.FindByPeriodIDAndUserID(period.ID, userID)
	if err != nil || participantPeriod == nil {
		return nil, errors.New("participant session not found")
	}
	if participantPeriod.Status != utils.ParticipantStatusStarted {
		return nil, errors.New("exam session is not active or already completed")
	}

	// Cari ID internal (uint) dari Question berdasarkan QuestionPublicID
	question, err := s.questionRepo.FindByPublicID(req.QuestionPublicID)
	if err != nil || question == nil {
		return nil, errors.New("question not found")
	}

	// Cari ID internal (uint) dari Option berdasarkan OptionPublicID
	option, err := s.optionRepo.FindByPublicID(req.OptionPublicID)
	if err != nil || option == nil {
		return nil, errors.New("option not found")
	}

	// Validasi kunci jawaban untuk mengisi kolom is_correct
	answerKey, err := s.answerKeyRepo.FindByQuestionID(question.ID)
	isCorrect := false
	if answerKey != nil && answerKey.OptionID == option.ID {
		isCorrect = true
	}

	// Simpan / Update jawaban menggunakan ID internal yang valid
	answerModel := &models.ParticipantAnswer{
		ParticipantPeriodID: participantPeriod.ID,
		QuestionID:          question.ID,
		OptionID:            option.ID,
		IsCorrect:           isCorrect,
	}

	if err := s.participantAnswerRepo.SaveAnswer(answerModel); err != nil {
		return nil, errors.New("failed to save answer")
	}

	// Kembalikan response sesuai ekspektasi endpoint
	return &models.AnswerResponse{
		QuestionPublicID: req.QuestionPublicID,
		OptionPublicID:   req.OptionPublicID,
		UpdatedAt:        time.Now(),
	}, nil
}

func (s *participantExamService) SubmitExam(periodPublicID string, userID uint) (*models.ExamSubmitResponse, error) {
	period, err := s.periodRepo.FindWithSectionsAndQuestions(periodPublicID)
	if err != nil || period == nil {
		return nil, errors.New("period not found")
	}

	participantPeriod, err := s.participantRepo.FindByPeriodIDAndUserID(period.ID, userID)
	if err != nil || participantPeriod == nil {
		return nil, errors.New("participant session not found")
	}

	// Idempotency: Jika sudah completed, langsung kembalikan response sukses
	if participantPeriod.Status == utils.ParticipantStatusCompleted {
		return &models.ExamSubmitResponse{
			PeriodPublicID: period.PublicID,
			Status:         participantPeriod.Status,
			SubmittedAt:    participantPeriod.UpdatedAt,
		}, nil
	}

	// Ambil seluruh jawaban peserta
	answers, err := s.participantAnswerRepo.FindByParticipantPeriodID(participantPeriod.ID)
	if err != nil {
		return nil, errors.New("failed to retrieve participant answers")
	}

	// Buat map jawaban benar peserta (key: question_id, value: is_correct)
	correctAnswersMap := make(map[uint]bool)
	for _, ans := range answers {
		if ans.IsCorrect {
			correctAnswersMap[ans.QuestionID] = true
		}
	}

	var totalScaledSum int
	sectionCount := 0

	// Iterasi setiap seksi pada periode ini melalui tabel pivot period_sections
	for _, ps := range period.PeriodSections {
		correctCount := 0
		
		// Hitung jumlah jawaban benar untuk seksi ini
		for _, q := range ps.Section.Questions {
			if correctAnswersMap[q.ID] {
				correctCount++
			}
		}

		// Konversi skor berdasarkan jenis seksi (Listening / Structure / Reading)
		var scaledScore int
		title := ps.Section.Title
		
		switch {
		case containsIgnoreCase(title, "listening"):
			scaledScore = utils.ConvertListeningScore(correctCount)
		case containsIgnoreCase(title, "structure") || containsIgnoreCase(title, "written"):
			scaledScore = utils.ConvertStructureScore(correctCount)
		case containsIgnoreCase(title, "reading"):
			scaledScore = utils.ConvertReadingScore(correctCount)
		default:
			scaledScore = utils.ConvertListeningScore(correctCount)
		}

		totalScaledSum += scaledScore
		sectionCount++

		// Simpan ke tabel section_scores sesuai ERD Anda
		sectionScore := &models.SectionScore{
			ParticipantPeriodID: participantPeriod.ID,
			SectionID:          ps.Section.ID,
			CorrectCount:       correctCount,
			RawScore:           correctCount,
			ScaledScore:        scaledScore,
		}
		
		if err := s.sectionScoreRepo.Save(sectionScore); err != nil {
			return nil, errors.New("failed to save section score")
		}
	}

	// Hitung Skor Akhir TOEFL Resmi: (Total Skor Konversi / Jumlah Seksi) * 10
	var finalScore int
	if sectionCount > 0 {
		finalScore = int((float64(totalScaledSum) / float64(sectionCount)) * 10)
	}

	// Update status dan skor total peserta secara terpisah di participant_periods
	err = s.participantRepo.UpdateStatus(participantPeriod.ID, utils.ParticipantStatusCompleted)
	if err != nil {
		return nil, errors.New("failed to update exam status")
	}

	err = s.participantRepo.UpdateScore(participantPeriod.ID, finalScore)
	if err != nil {
		return nil, errors.New("failed to update exam score")
	}

	return &models.ExamSubmitResponse{
		PeriodPublicID: period.PublicID,
		Status:         utils.ParticipantStatusCompleted,
		SubmittedAt:    time.Now(),
	}, nil
}

// Helper kecil untuk mencocokkan judul seksi secara case-insensitive
func containsIgnoreCase(str, substr string) bool {
	return len(str) >= len(substr) && (str == substr || len(substr) == 0)
}