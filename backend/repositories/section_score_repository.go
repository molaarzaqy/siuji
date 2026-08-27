package repositories

import (
	"siuji-backend/models"

	"gorm.io/gorm"
)

type SectionScoreRepository interface {
	Save(sectionScore *models.SectionScore) error
}

type sectionScoreRepository struct {
	db *gorm.DB
}

func NewSectionScoreRepository(db *gorm.DB) SectionScoreRepository {
	return &sectionScoreRepository{
		db: db,
	}
}

func (r *sectionScoreRepository) Save(sectionScore *models.SectionScore) error {
	// Menggunakan Upsert sederhana berdasarkan participant_period_id dan section_id 
	// agar jika di-submit ulang tidak terjadi duplikasi data.
	var existing models.SectionScore
	err := r.db.Where("participant_period_id = ? AND section_id = ?", sectionScore.ParticipantPeriodID, sectionScore.SectionID).First(&existing).Error

	if err == nil {
		// Jika sudah ada, update
		existing.CorrectCount = sectionScore.CorrectCount
		existing.RawScore = sectionScore.RawScore
		existing.ScaledScore = sectionScore.ScaledScore
		return r.db.Save(&existing).Error
	} else if err == gorm.ErrRecordNotFound {
		// Jika belum ada, buat baru
		return r.db.Create(sectionScore).Error
	}

	return err
}