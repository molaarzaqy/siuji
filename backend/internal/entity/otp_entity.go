package entity

import (
	"time"

	"gorm.io/gorm"
)

type OTP struct {
	ID        uint           `gorm:"column:id;primaryKey"`
	Email     string         `gorm:"column:email"`
	Code      string         `gorm:"column:code"`
	Purpose   string         `gorm:"column:purpose"`
	ExpiresAt time.Time      `gorm:"column:expires_at"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (OTP) TableName() string {
	return "otps"
}

func (o *OTP) IsExpired() bool {
	return time.Now().After(o.ExpiresAt)
}