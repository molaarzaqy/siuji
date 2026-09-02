package repository

import (
	"errors"

	"siuji-backend/internal/entity"

	"gorm.io/gorm"
)

type OptionRepository interface {
	Create(option *entity.Option) error
	FindByPublicID(publicID string) (*entity.Option, error)
	GetMaxPositionInQuestion(questionID uint) (int, error)
	Update(option *entity.Option) error
	Delete(publicID string) error
	UpdatePositionsByPublicIDs(publicIDs []string) error
}

type optionRepository struct {
	db *gorm.DB
}

func NewOptionRepository(db *gorm.DB) OptionRepository {
	return &optionRepository{db: db}
}

func (r *optionRepository) Create(option *entity.Option) error {
	return r.db.Create(option).Error
}

func (r *optionRepository) FindByPublicID(publicID string) (*entity.Option, error) {
	var option entity.Option
	err := r.db.Where("public_id = ?", publicID).First(&option).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("option not found")
		}
		return nil, err
	}
	return &option, nil
}

func (r *optionRepository) GetMaxPositionInQuestion(questionID uint) (int, error) {
	var maxPosition int
	err := r.db.Model(&entity.Option{}).
		Where("question_id = ?", questionID).
		Select("COALESCE(MAX(position), 0)").
		Scan(&maxPosition).Error
	return maxPosition, err
}

func (r *optionRepository) Update(option *entity.Option) error {
	return r.db.Save(option).Error
}

func (r *optionRepository) Delete(publicID string) error {
	result := r.db.Where("public_id = ?", publicID).Delete(&entity.Option{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("option not found")
	}
	return nil
}

func (r *optionRepository) UpdatePositionsByPublicIDs(publicIDs []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for index, publicID := range publicIDs {
			result := tx.Model(&entity.Option{}).
				Where("public_id = ?", publicID).
				Update("position", index+1)
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return errors.New("one or more options not found")
			}
		}
		return nil
	})
}