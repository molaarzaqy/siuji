package repositories

import (
	"errors"
	"siuji-backend/models"

	"gorm.io/gorm"
)

type AnswerKeyRepository interface {
	FindByQuestionID(questionID uint) (*models.AnswerKey, error)
}

type answerKeyRepository struct {
	db *gorm.DB
}

func NewAnswerKeyRepository(db *gorm.DB) AnswerKeyRepository {
	return &answerKeyRepository{
		db: db,
	}
}

func (r *answerKeyRepository) FindByQuestionID(questionID uint) (*models.AnswerKey, error) {
	var answerKey models.AnswerKey
	err := r.db.Where("question_id = ?", questionID).First(&answerKey).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &answerKey, nil
}