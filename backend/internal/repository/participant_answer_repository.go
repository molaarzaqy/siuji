package repository

import (
	"errors"
	"time"

	"siuji-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ParticipantAnswerRepository interface {
	Upsert(answer *entity.ParticipantAnswer) error
	FindByParticipantPeriodAndQuestion(ppID, questionID uint) (*entity.ParticipantAnswer, error)
	FindAllByParticipantPeriod(ppID uint) ([]entity.ParticipantAnswer, error)
	CountCorrectByParticipantPeriodAndSection(ppID, sectionID uint) (int, error)
}

type participantAnswerRepository struct {
	db *gorm.DB
}

func NewParticipantAnswerRepository(db *gorm.DB) ParticipantAnswerRepository {
	return &participantAnswerRepository{db: db}
}

func (r *participantAnswerRepository) Upsert(answer *entity.ParticipantAnswer) error {
	if answer.PublicID == uuid.Nil {
		answer.PublicID = uuid.New()
	}
	now := time.Now()
	if answer.CreatedAt.IsZero() {
		answer.CreatedAt = now
	}
	answer.UpdatedAt = now

	return r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "participant_period_id"},
			{Name: "question_id"},
		},
		DoUpdates: clause.AssignmentColumns([]string{"option_id", "is_correct", "updated_at"}),
	}).Create(answer).Error
}

func (r *participantAnswerRepository) FindByParticipantPeriodAndQuestion(ppID, questionID uint) (*entity.ParticipantAnswer, error) {
	var answer entity.ParticipantAnswer
	err := r.db.
		Where("participant_period_id = ? AND question_id = ?", ppID, questionID).
		First(&answer).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &answer, nil
}

func (r *participantAnswerRepository) FindAllByParticipantPeriod(ppID uint) ([]entity.ParticipantAnswer, error) {
	var list []entity.ParticipantAnswer
	err := r.db.
		Where("participant_period_id = ?", ppID).
		Find(&list).Error
	return list, err
}

func (r *participantAnswerRepository) CountCorrectByParticipantPeriodAndSection(ppID, sectionID uint) (int, error) {
	var count int64
	err := r.db.Model(&entity.ParticipantAnswer{}).
		Joins("JOIN questions ON questions.id = participant_answers.question_id").
		Where("participant_answers.participant_period_id = ? AND questions.section_id = ? AND participant_answers.is_correct = true", ppID, sectionID).
		Count(&count).Error
	return int(count), err
}

