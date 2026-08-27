package repositories

import (
	"errors"
	"siuji-backend/models"

	"gorm.io/gorm"
)

type ParticipantAnswerRepository interface {
	SaveAnswer(answer *models.ParticipantAnswer) error
	FindByParticipantAndQuestion(participantPeriodID uint, questionID uint) (*models.ParticipantAnswer, error)
	FindByParticipantPeriodID(participantPeriodID uint) ([]models.ParticipantAnswer, error)
}

type participantAnswerRepository struct {
	db *gorm.DB
}

func NewParticipantAnswerRepository(db *gorm.DB) ParticipantAnswerRepository {
	return &participantAnswerRepository{
		db: db,
	}
}

func (r *participantAnswerRepository) SaveAnswer(answer *models.ParticipantAnswer) error {
	var existing models.ParticipantAnswer
	err := r.db.Where("participant_period_id = ? AND question_id = ?", answer.ParticipantPeriodID, answer.QuestionID).First(&existing).Error
	
	if err == nil {
		// Jika sudah ada, update option_id dan is_correct-nya
		existing.OptionID = answer.OptionID
		existing.IsCorrect = answer.IsCorrect
		return r.db.Save(&existing).Error
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// Jika belum ada, buat baru (Insert)
		return r.db.Create(answer).Error
	}
	
	return err
}

func (r *participantAnswerRepository) FindByParticipantAndQuestion(participantPeriodID uint, questionID uint) (*models.ParticipantAnswer, error) {
	var answer models.ParticipantAnswer
	err := r.db.Where("participant_period_id = ? AND question_id = ?", participantPeriodID, questionID).First(&answer).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &answer, nil
}

func (r *participantAnswerRepository) FindByParticipantPeriodID(participantPeriodID uint) ([]models.ParticipantAnswer, error) {
    var answers []models.ParticipantAnswer
    err := r.db.Where("participant_period_id = ?", participantPeriodID).Find(&answers).Error
    return answers, err
}