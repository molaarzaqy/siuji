package repository

import (
	"errors"

	"siuji-backend/internal/entity"

	"gorm.io/gorm"
)

type PeriodSectionRepository interface {
	Create(ps *entity.PeriodSection) error
	FindByPublicID(publicID string) (*entity.PeriodSection, error)
	FindByPeriodID(periodID uint) ([]entity.PeriodSection, error)
	ExistsByPeriodAndSection(periodID, sectionID uint) (bool, error)
	UpdatePosition(publicID string, position int) error
	CountByPeriodID(periodID uint) (int64, error)
	UpdatePositionsBulk(periodID uint, positions map[uint]int) error
	DeleteByPeriodAndSection(periodID, sectionID uint) error
}

type periodSectionRepository struct {
	db *gorm.DB
}

func NewPeriodSectionRepository(db *gorm.DB) PeriodSectionRepository {
	return &periodSectionRepository{db: db}
}

func (r *periodSectionRepository) Create(ps *entity.PeriodSection) error {
	return r.db.Create(ps).Error
}

func (r *periodSectionRepository) FindByPublicID(publicID string) (*entity.PeriodSection, error) {
	var ps entity.PeriodSection
	err := r.db.Preload("Section").Preload("Period").
		Where("public_id = ?", publicID).First(&ps).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("period section not found")
		}
		return nil, err
	}
	return &ps, nil
}

func (r *periodSectionRepository) FindByPeriodID(periodID uint) ([]entity.PeriodSection, error) {
	var list []entity.PeriodSection
	err := r.db.Preload("Section").
		Where("period_id = ?", periodID).
		Order("position ASC").
		Find(&list).Error
	return list, err
}

func (r *periodSectionRepository) ExistsByPeriodAndSection(periodID, sectionID uint) (bool, error) {
	var count int64
	err := r.db.Model(&entity.PeriodSection{}).
		Where("period_id = ? AND section_id = ?", periodID, sectionID).
		Count(&count).Error
	return count > 0, err
}

func (r *periodSectionRepository) UpdatePosition(publicID string, position int) error {
	result := r.db.Model(&entity.PeriodSection{}).
		Where("public_id = ?", publicID).
		Update("position", position)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("period section not found")
	}
	return nil
}

func (r *periodSectionRepository) CountByPeriodID(periodID uint) (int64, error) {
	var count int64
	err := r.db.Model(&entity.PeriodSection{}).
		Where("period_id = ?", periodID).
		Count(&count).Error
	return count, err
}

func (r *periodSectionRepository) UpdatePositionsBulk(periodID uint, positions map[uint]int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for sectionID, position := range positions {
			result := tx.Model(&entity.PeriodSection{}).
				Where("period_id = ? AND section_id = ?", periodID, sectionID).
				Update("position", position)
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return errors.New("one or more sections not found in this period")
			}
		}
		return nil
	})
}

func (r *periodSectionRepository) DeleteByPeriodAndSection(periodID, sectionID uint) error {
	result := r.db.Where("period_id = ? AND section_id = ?", periodID, sectionID).
		Delete(&entity.PeriodSection{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("section not found in this period")
	}
	return nil
}