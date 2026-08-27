package repositories

import (
	"errors"
	"siuji-backend/models"

	"gorm.io/gorm"
)

type PeriodRepository interface {
	Create(period *models.Period) error
	FindAllPagination(filter, sort string, limit, offset int) ([]models.Period, int64, error)
	FindByPublicID(publicID string) (*models.Period, error)
	Update(period *models.Period) error
	Delete(publicID string) error
	UpdateStatus(id uint, status string) error
	// period - section management (pivot & relations)
	AddSectionToPeriod(periodSection *models.PeriodSection) error
	RemoveSectionFromPeriod(periodID uint, sectionID uint) error
	GetMaxPositionInPeriod(periodID uint) (int, error)
	UpdateSectionPositions(updates []SectionPositionUpdate) error
	FindPeriodSection(periodID uint, sectionID uint) (*models.PeriodSection, error)
	FindWithSectionsAndQuestions(publicID string) (*models.Period, error)
}

type SectionPositionUpdate struct {
	PeriodID  uint
	SectionID uint
	Position  int
}

type periodRepository struct {
	db *gorm.DB
}

func NewPeriodRepository(db *gorm.DB) PeriodRepository {
	return &periodRepository{db: db}
}

func (r *periodRepository) Create(period *models.Period) error {
	return r.db.Create(period).Error
}

func (r *periodRepository) FindAllPagination(filter, sort string, limit, offset int) ([]models.Period, int64, error) {
	var periods []models.Period
	var totalData int64

	query := r.db.Model(&models.Period{})
	
	// Filter sederhana jika ada
	if filter != "" {
		query = query.Where("title ILIKE ? OR status ILIKE ?", "%"+filter+"%", "%"+filter+"%")
	}
	// Hitung total data
	if err := query.Count(&totalData).Error; err != nil {
		return nil, 0, err
	}
	// Sorting
	if sort != "" {
		query = query.Order(sort)
	} else {
		query = query.Order("created_at DESC")
	}
	// Paginasi & Preload
	err := query.Limit(limit).Offset(offset).Preload("PeriodSections.Section").Find(&periods).Error
	return periods, totalData, err
}

func (r *periodRepository) FindByPublicID(publicID string) (*models.Period, error) {
	var period models.Period
	err := r.db.Preload("PeriodSections.Section").Where("public_id = ? ", publicID).First(&period).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &period, err
}

func (r *periodRepository) Update(period *models.Period) error {
	return r.db.Save(period).Error
}

func (r *periodRepository) Delete(publicID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var period models.Period
		if err := tx.Where("public_id = ?", publicID).First(&period).Error; err != nil {
			return err
		}
		if err := tx.Where("period_id = ?", period.ID).Delete(&models.PeriodSection{}).Error; err != nil {
			return err
		}
		return tx.Delete(&period).Error
	})
}

func (r *periodRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&models.Period{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *periodRepository) AddSectionToPeriod(periodSection *models.PeriodSection) error {
	return r.db.Create(periodSection).Error
}

func (r *periodRepository) RemoveSectionFromPeriod(periodID uint, sectionID uint) error {
	return r.db.Where("period_id = ? AND section_id = ?", periodID, sectionID).Delete(&models.PeriodSection{}).Error
}

func (r *periodRepository) GetMaxPositionInPeriod(periodID uint) (int, error) {
	var maxPosition int
	err := r.db.Model(&models.PeriodSection{}).
			Where("period_id = ?", periodID).Select("COALESCE(MAX(position), 0)").
			Scan(&maxPosition).Error
	return maxPosition, err
}

func (r *periodRepository) UpdateSectionPositions(updates []SectionPositionUpdate) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, update := range updates {
			err := tx.Model(&models.PeriodSection{}).
					Where("period_id = ? AND section_id = ?", update.PeriodID, update.SectionID).
					Update("position", update.Position).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *periodRepository) FindPeriodSection(periodID uint, sectionID uint) (*models.PeriodSection, error) {
	var periodSection models.PeriodSection
	err := r.db.Where("period_id = ? AND section_id = ?", periodID, sectionID).First(&periodSection).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &periodSection, nil
}

func (r *periodRepository) FindWithSectionsAndQuestions(publicID string) (*models.Period, error) {
    var period models.Period
    err := r.db.
        Preload("PeriodSections.Section.Questions.Options").
        Where("public_id = ?", publicID).
        First(&period).Error
    
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &period, nil
}