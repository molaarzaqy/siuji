package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID         uint 		  `gorm:"column:id;primaryKey"`
	PublicID   uuid.UUID	  `gorm:"column:public_id"`
	Name       string         `gorm:"column:name"`
	Email      string         `gorm:"column:email"`
	University string         `gorm:"column:university"`
	NIM        string         `gorm:"column:nim"`
	Password   string         `gorm:"column:password"`
	Role       string         `gorm:"column:role"`
	CreatedAt  time.Time      `gorm:"column:created_at"`
	UpdatedAt  time.Time      `gorm:"column:updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;index"`

	// relations
	ParticipantPeriods []ParticipantPeriod `gorm:"foreignKey:UserID"`
}

func (User) TableName() string {
	return "users"
}