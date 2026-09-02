package repository

import (
	"errors"

	"siuji-backend/internal/entity"

	"gorm.io/gorm"
)

type QuestionRepository interface {
	Create(question *entity.Question) error
	FindByPublicID(publicID string) (*entity.Question, error)
	FindByPublicIDWithOptions(publicID string) (*entity.Question, error)
	GetMaxPositionInSection(sectionID uint) (int, error)
	Update(question *entity.Question) error
	Delete(publicID string) error
	UpdatePositionsByPublicIDs(publicIDs []string) error
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
	var question entity.Question
	err := r.db.Where("public_id = ?", publicID).First(&question).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("question not found")
		}
		return nil, err
	}
	return &question, nil
}

func (r *questionRepository) FindByPublicIDWithOptions(publicID string) (*entity.Question, error) {
	var question entity.Question
	err := r.db.
		Preload("Options", func(db *gorm.DB) *gorm.DB {
			return db.Order("options.position ASC")
		}).
		Preload("AnswerKeys.Option").
		Where("public_id = ?", publicID).
		First(&question).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("question not found")
		}
		return nil, err
	}
	return &question, nil
}

func (r *questionRepository) GetMaxPositionInSection(sectionID uint) (int, error) {
	var maxPosition int
	err := r.db.Model(&entity.Question{}).
		Where("section_id = ?", sectionID).
		Select("COALESCE(MAX(position), 0)").
		Scan(&maxPosition).Error
	return maxPosition, err
}

func (r *questionRepository) Update(question *entity.Question) error {
	return r.db.Save(question).Error
}

func (r *questionRepository) Delete(publicID string) error {
	result := r.db.Where("public_id = ?", publicID).Delete(&entity.Question{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("question not found")
	}
	return nil
}

func (r *questionRepository) UpdatePositionsByPublicIDs(publicIDs []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for index, publicID := range publicIDs {
			result := tx.Model(&entity.Question{}).
				Where("public_id = ?", publicID).
				Update("position", index+1)
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