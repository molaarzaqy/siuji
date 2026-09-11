package repository

import (
	"errors"
	"siuji-backend/internal/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ScoreConversionRepository interface {
	FindScaledScore(sectionType string, correctCount int) (int, error)
	BulkUpsert(list []entity.ScoreConversion) error
	Create(sc *entity.ScoreConversion) error
	FindAll(sectionType string) ([]entity.ScoreConversion, error)
	FindByID(id uint) (*entity.ScoreConversion, error)
	Update(sc *entity.ScoreConversion) error
	Delete(id uint) error
}

type scoreConversionRepository struct {
	db *gorm.DB
}

func NewScoreConversionRepository(db *gorm.DB) ScoreConversionRepository {
	return &scoreConversionRepository{
		db: db,
	}
}

// FindScaledScore mencari scaled_score untuk section_type & correct_count tertentu.
// Database telah di-seed dengan data standar TOEFL ITP/PBT (0-50 Listening, 0-40 Structure, 0-50 Reading).
// Penggunaan `correct_count <= ?` dan order DESC berfungsi sebagai fallback aman bila ada nilai di luar batas.
func (r *scoreConversionRepository) FindScaledScore(sectionType string, correctCount int) (int, error) {
	var conversion entity.ScoreConversion
	err := r.db.Where("section_type = ? AND correct_count <= ?", sectionType, correctCount).
		Order("correct_count DESC").
		First(&conversion).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errors.New("no score conversion found for this section type")
		}
		return 0, err
	}
	return conversion.ScaledScore, nil
}

func (r *scoreConversionRepository) BulkUpsert(list []entity.ScoreConversion) error {
	if len(list) == 0 {
		return nil
	}
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "section_type"}, {Name: "correct_count"}},
		DoUpdates: clause.AssignmentColumns([]string{"scaled_score"}),
	}).Create(&list).Error
}

func (r *scoreConversionRepository) Create(sc *entity.ScoreConversion) error {
	return r.db.Create(sc).Error
}

func (r *scoreConversionRepository) FindAll(sectionType string) ([]entity.ScoreConversion, error) {
	var list []entity.ScoreConversion
	db := r.db.Model(&entity.ScoreConversion{})
	if sectionType != "" {
		db = db.Where("section_type = ?", sectionType)
	}
	err := db.Order("section_type ASC, correct_count ASC").Find(&list).Error
	return list, err
}

func (r *scoreConversionRepository) FindByID(id uint) (*entity.ScoreConversion, error) {
	var sc entity.ScoreConversion
	err := r.db.First(&sc, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("score conversion not found")
		}
		return nil, err
	}
	return &sc, nil
}

func (r *scoreConversionRepository) Update(sc *entity.ScoreConversion) error {
	return r.db.Save(sc).Error
}

func (r *scoreConversionRepository) Delete(id uint) error {
	result := r.db.Delete(&entity.ScoreConversion{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("score conversion not found")
	}
	return nil
}