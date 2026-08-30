package repository

import (
	"errors"
	"siuji-backend/internal/entity"
	"time"

	"gorm.io/gorm"
)

type OTPRepository interface {
	Create(otp *entity.OTP) error
	FindValidByEmailAndCode(email, code, purpose string) (*entity.OTP, error)
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

func (r *otpRepository) FindValidByEmailAndCode(email, code, purpose string) (*entity.OTP, error) {
	var otp entity.OTP
	now := time.Now()

	err := r.db.
		Where("email = ? AND code = ? AND purpose = ? AND expires_at > ?", email, code, purpose, now).
		First(&otp).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid or expired OTP")
		}
		return nil, err
	}
	return &otp, nil
}

func (r *otpRepository) DeleteByEmail(email string) error {
	return r.db.Where("email = ?", email).Delete(&entity.OTP{}).Error
}

func (r *otpRepository) DeleteExpired() error {
	now := time.Now()
	return r.db.Where("expires_at < ?", now).Delete(&entity.OTP{}).Error
}