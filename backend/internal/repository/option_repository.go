package repository

import (
	"errors"

	"siuji-backend/internal/entity"

	"github.com/google/uuid"
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
	parsedID, err := uuid.Parse(publicID)
	if err != nil {
		return nil, errors.New("invalid uuid format")
	}

	var option entity.Option
	err = r.db.Where("public_id = ?", parsedID).First(&option).Error
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
	parsedID, err := uuid.Parse(publicID)
	if err != nil {
		return errors.New("invalid uuid format")
	}

	result := r.db.Where("public_id = ?", parsedID).Delete(&entity.Option{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("option not found")
	}
	return nil
}

func (r *optionRepository) UpdatePositionsByPublicIDs(publicIDs []string) error {
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
			position := index + 1
			result := tx.Model(&entity.Option{}).
				Where("public_id = ?", parsedID).
				Updates(map[string]interface{}{
					"position": position,
					"label":    entity.LabelFromPosition(position),
				})
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