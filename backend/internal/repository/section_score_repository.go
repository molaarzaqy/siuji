package repository

import (
	"time"

	"siuji-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SectionScoreRepository interface {
	Upsert(score *entity.SectionScore) error
	FindByParticipantPeriodID(ppID uint) ([]entity.SectionScore, error)
}

type sectionScoreRepository struct {
	db *gorm.DB
}

func NewSectionScoreRepository(db *gorm.DB) SectionScoreRepository {
	return &sectionScoreRepository{db: db}
}

func (r *sectionScoreRepository) Upsert(score *entity.SectionScore) error {
	if score.PublicID == uuid.Nil {
		score.PublicID = uuid.New()
	}
	now := time.Now()
	if score.CreatedAt.IsZero() {
		score.CreatedAt = now
	}
	score.UpdatedAt = now

	return r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "participant_period_id"},
			{Name: "section_id"},
		},
		DoUpdates: clause.AssignmentColumns([]string{"correct_count", "raw_score", "scaled_score", "updated_at"}),
	}).Create(score).Error
}

func (r *sectionScoreRepository) FindByParticipantPeriodID(ppID uint) ([]entity.SectionScore, error) {
	var list []entity.SectionScore
	err := r.db.
		Preload("Section").
		Where("participant_period_id = ?", ppID).
		Joins("JOIN period_sections ON period_sections.section_id = section_scores.section_id").
		Joins("JOIN participant_periods ON participant_periods.id = section_scores.participant_period_id AND participant_periods.period_id = period_sections.period_id").
		Order("period_sections.position ASC").
		Find(&list).Error
	if err != nil {
		// Fallback without ordering by period_sections if join fails
		err = r.db.
			Preload("Section").
			Where("participant_period_id = ?", ppID).
			Find(&list).Error
	}
	return list, err
}

