package repository

import (
	"errors"

	"siuji-backend/internal/entity"

	"gorm.io/gorm"
)

type SectionRepository interface {
	Create(section *entity.Section) error
	FindByPublicID(publicID string) (*entity.Section, error)
	FindAll() ([]entity.Section, error)
	Update(section *entity.Section) error
	Delete(publicID string) error
}

type sectionRepository struct {
	db *gorm.DB
}

func NewSectionRepository(db *gorm.DB) SectionRepository {
	return &sectionRepository{db: db}
}

func (r *sectionRepository) Create(section *entity.Section) error {
	return r.db.Create(section).Error
}

func (r *sectionRepository) FindByPublicID(publicID string) (*entity.Section, error) {
	var section entity.Section
	err := r.db.Where("public_id = ?", publicID).First(&section).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("section not found")
		}
		return nil, err
	}
	return &section, nil
}

func (r *sectionRepository) FindAll() ([]entity.Section, error) {
	var sections []entity.Section
	err := r.db.Order("title ASC").Find(&sections).Error
	return sections, err
}

func (r *sectionRepository) Update(section *entity.Section) error {
	return r.db.Save(section).Error
}

func (r *sectionRepository) Delete(publicID string) error {
	result := r.db.Where("public_id = ?", publicID).Delete(&entity.Section{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("section not found")
	}
	return nil
}