package entity

import (
	"time"

	"github.com/google/uuid"
)

type ParticipantPeriod struct {
	ID        uint      `gorm:"column:id;primaryKey"`
	PublicID  uuid.UUID `gorm:"column:public_id"`
	UserID    uint      `gorm:"column:user_id"`
	PeriodID  uint      `gorm:"column:period_id"`
	Status    string    `gorm:"column:status"`
	Score     *int      `gorm:"column:score"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`

	// Relations
	User   User   `gorm:"foreignKey:UserID"`
	Period Period `gorm:"foreignKey:PeriodID"`
}

func (ParticipantPeriod) TableName() string {
	return "participant_periods"
}