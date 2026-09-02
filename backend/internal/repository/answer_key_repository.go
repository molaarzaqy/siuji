package repository

import (
	"errors"

	"siuji-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AnswerKeyRepository interface {
	FindByQuestionID(questionID uint) (*entity.AnswerKey, error)
	Upsert(questionID, optionID uint) (*entity.AnswerKey, error)
}

type answerKeyRepository struct {
	db *gorm.DB
}

func NewAnswerKeyRepository(db *gorm.DB) AnswerKeyRepository {
	return &answerKeyRepository{db: db}
}

func (r *answerKeyRepository) FindByQuestionID(questionID uint) (*entity.AnswerKey, error) {
	var answerKey entity.AnswerKey
	err := r.db.Preload("Option").Where("question_id = ?", questionID).First(&answerKey).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &answerKey, nil
}

func (r *answerKeyRepository) Upsert(questionID, optionID uint) (*entity.AnswerKey, error) {
	var answerKey entity.AnswerKey

	err := r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("question_id = ?", questionID).First(&answerKey).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				answerKey = entity.AnswerKey{
					PublicID:   uuid.New(),
					QuestionID: questionID,
					OptionID:   optionID,
				}
				return tx.Create(&answerKey).Error
			}
			return err
		}
		answerKey.OptionID = optionID
		return tx.Save(&answerKey).Error
	})
	if err != nil {
		return nil, err
	}

	if err := r.db.Preload("Option").First(&answerKey, answerKey.ID).Error; err != nil {
		return nil, err
	}
	return &answerKey, nil
}