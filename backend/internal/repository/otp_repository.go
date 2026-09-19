package repository

import (
	"errors"
	"siuji-backend/internal/entity"
	"time"

	"gorm.io/gorm"
)

type OTPRepository interface {
	Create(otp *entity.OTP) error
	FindValidByEmailAndPurpose(email, purpose string) (*entity.OTP, error)
	IncrementAttempts(id uint) error
	DeleteByEmail(email string) error
	DeleteExpired() error
}

type otpRepository struct {
	db *gorm.DB
}

func NewOTPRepository(db *gorm.DB) OTPRepository {
	return &otpRepository{db: db}
}

func (r *otpRepository) Create(otp *entity.OTP) error {
	return r.db.Create(otp).Error
}

func (r *otpRepository) FindValidByEmailAndPurpose(email, purpose string) (*entity.OTP, error) {
	var otp entity.OTP
	err := r.db.
		Where("email = ? AND purpose = ? AND expires_at > ?", email, purpose, time.Now()).
		Order("created_at DESC").
		First(&otp).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("otp not found")
		}
		return nil, err
	}
	return &otp, nil
}

func (r *otpRepository) IncrementAttempts(id uint) error {
	return r.db.Model(&entity.OTP{}).
		Where("id = ?", id).
		UpdateColumn("attempts", gorm.Expr("attempts + 1")).Error
}

func (r *otpRepository) DeleteByEmail(email string) error {
	return r.db.Where("email = ?", email).Delete(&entity.OTP{}).Error
}

func (r *otpRepository) DeleteExpired() error {
	now := time.Now()
	return r.db.Where("expires_at < ?", now).Delete(&entity.OTP{}).Error
}