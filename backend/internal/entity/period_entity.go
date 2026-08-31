package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Period struct {
	ID                  uint           `gorm:"column:id;primaryKey"`
	PublicID            uuid.UUID      `gorm:"column:public_id"`
	Title               string         `gorm:"column:title"`
	Month               string         `gorm:"column:month"`
	Year                int            `gorm:"column:year"`
	Status              string         `gorm:"column:status"`
	CertificateURL      string         `gorm:"column:certificate_url"`
	CertificateExpMonth time.Time      `gorm:"column:certificate_exp_month"`
	MinPassingGrade     int            `gorm:"column:min_passing_grade"`
	MaxPassingGrade     int            `gorm:"column:max_passing_grade"`
	StartTime           time.Time      `gorm:"column:start_time"`
	EndTime             time.Time      `gorm:"column:end_time"`
	CreatedAt           time.Time      `gorm:"column:created_at"`
	UpdatedAt           time.Time      `gorm:"column:updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"column:deleted_at;index"`

	// Relations
	PeriodSections []PeriodSection `gorm:"foreignKey:PeriodID"`
	ParticipantPeriods []ParticipantPeriod `gorm:"foreignKey:PeriodID"`
}

func (Period) TableName() string {
	return "periods"
}