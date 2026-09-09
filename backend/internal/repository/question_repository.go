package repository

import (
	"errors"

	"siuji-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuestionRepository interface {
	Create(question *entity.Question) error
	FindByPublicID(publicID string) (*entity.Question, error)
	FindByPublicIDWithOptions(publicID string) (*entity.Question, error)
	GetMaxNumberInSection(sectionID uint) (int, error)
	Update(question *entity.Question) error
	Delete(publicID string) error
	UpdateNumbersByPublicIDs(publicIDs []string) error
}

type questionRepository struct {
	db *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) QuestionRepository {
	return &questionRepository{db: db}
}

func (r *questionRepository) Create(question *entity.Question) error {
	return r.db.Create(question).Error
}

func (r *questionRepository) FindByPublicID(publicID string) (*entity.Question, error) {
	parsedID, err := uuid.Parse(publicID)
	if err != nil {
		return nil, errors.New("invalid uuid format")
	}

	var question entity.Question
	err = r.db.Where("public_id = ?", parsedID).First(&question).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("question not found")
		}
		return nil, err
	}
	return &question, nil
}

func (r *questionRepository) FindByPublicIDWithOptions(publicID string) (*entity.Question, error) {
	parsedID, err := uuid.Parse(publicID)
	if err != nil {
		return nil, errors.New("invalid uuid format")
	}

	var question entity.Question
	err = r.db.
		Preload("Options", func(db *gorm.DB) *gorm.DB {
			return db.Order("options.position ASC")
		}).
		Preload("AnswerKeys.Option").
		Where("public_id = ?", parsedID).
		First(&question).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("question not found")
		}
		return nil, err
	}
	return &question, nil
}

func (r *questionRepository) GetMaxNumberInSection(sectionID uint) (int, error) {
	var maxNumber int
	err := r.db.Model(&entity.Question{}).
		Where("section_id = ?", sectionID).
		Select("COALESCE(MAX(number), 0)").
		Scan(&maxNumber).Error
	return maxNumber, err
}

func (r *questionRepository) Update(question *entity.Question) error {
	return r.db.Save(question).Error
}

func (r *questionRepository) Delete(publicID string) error {
	parsedID, err := uuid.Parse(publicID)
	if err != nil {
		return errors.New("invalid uuid format")
	}

	result := r.db.Where("public_id = ?", parsedID).Delete(&entity.Question{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("question not found")
	}
	return nil
}

func (r *questionRepository) UpdateNumbersByPublicIDs(publicIDs []string) error {
	parsedIDs := make([]uuid.UUID, len(publicIDs))
	for i, idStr := range publicIDs {
		parsedID, err := uuid.Parse(idStr)
		if err != nil {
			return errors.New("invalid uuid format")
		}
		parsedIDs[i] = parsedID
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		for index, parsedID := range parsedIDs {
			result := tx.Model(&entity.Question{}).
				Where("public_id = ?", parsedID).
				Update("number", index+1)
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return errors.New("one or more questions not found")
			}
		}
		return nil
	})
}