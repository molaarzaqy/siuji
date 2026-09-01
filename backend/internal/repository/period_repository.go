package repository

import (
	"errors"
	"siuji-backend/internal/entity"

	"gorm.io/gorm"
)

type PeriodRepository interface {
	Create(period *entity.Period) error
	FindByPublicID(publicID string) (*entity.Period, error)
	FindByPublicIDWithSections(publicID string) (*entity.Period, error)
	FindAllPagination(filter, sort string, limit, offset int) ([]entity.Period, int64, error)
	Update(period *entity.Period) error
	Delete(publicID string) error
}

type periodRepository struct {
	db *gorm.DB
}

func NewPeriodRepository(db *gorm.DB) PeriodRepository {
	return &periodRepository{
		db: db,
	}
}

func (r *periodRepository) Create(period *entity.Period) error {
	return r.db.Create(period).Error
}

func (r *periodRepository) FindByPublicID(publicID string) (*entity.Period, error) {
	var period entity.Period
	err := r.db.Where("public_id = ?", publicID).First(&period).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("period not found")
		}
		return nil, err
	}
	return &period, nil
}

func (r *periodRepository) FindByPublicIDWithSections(publicID string) (*entity.Period, error) {
	var period entity.Period
	err := r.db.
		Preload("PeriodSections", func(db *gorm.DB) *gorm.DB {
			return db.Order("period_sections.position ASC")
		}).
		Preload("PeriodSections.Section").
		Where("public_id = ?", publicID).
		First(&period).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("period not found")
		}
		return nil, err
	}
	return &period, nil
}

func (r *periodRepository) FindAllPagination(filter, sort string, limit, offset int) ([]entity.Period, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	var periods []entity.Period
	var total int64

	db := r.db.Model(&entity.Period{})
	if filter != "" {
		db = db.Where("title ILIKE ?", "%"+filter+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	allowedSortFields := map[string]string{
		"title":       "title ASC",
		"-title":      "title DESC",
		"year":        "year ASC",
		"-year":       "year DESC",
		"status":      "status ASC",
		"-status":     "status DESC",
		"created_at":  "created_at ASC",
		"-created_at": "created_at DESC",
	}

	if sortClause, ok := allowedSortFields[sort]; ok {
		db = db.Order(sortClause)
	} else {
		db = db.Order("created_at DESC")
	}

	if err := db.Limit(limit).Offset(offset).Find(&periods).Error; err != nil {
		return nil, 0, err
	}
	return periods, total, nil
}

func (r *periodRepository) Update(period *entity.Period) error {
	return r.db.Save(period).Error
}

func (r *periodRepository) Delete(publicID string) error {
		result := r.db.Where("public_id = ?", publicID).Delete(&entity.Period{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("period not found")
	}
	return nil
}
