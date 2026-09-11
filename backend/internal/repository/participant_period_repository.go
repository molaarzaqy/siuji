package repository

import (
	"errors"

	"siuji-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ParticipantPeriodRepository interface {
	Create(pp *entity.ParticipantPeriod) error
	ExistsByPeriodAndUser(periodID, userID uint) (bool, error)
	FindByPeriodAndUserPublicID(periodID uint, userPublicID string) (*entity.ParticipantPeriod, error)
	FindAllByPeriodIDPagination(periodID uint, filter, sort string, limit, offset int) ([]entity.ParticipantPeriod, int64, error)
	Update(pp *entity.ParticipantPeriod) error
	DeleteByPeriodAndUserPublicID(periodID uint, userPublicID string) error
	BulkCreate(list []entity.ParticipantPeriod) error
	FindByUserID(userID uint) ([]entity.ParticipantPeriod, error)
	FindByPeriodIDAndUserID(periodID, userID uint) (*entity.ParticipantPeriod, error)
	FindByPeriodPublicIDAndUserID(periodPublicID string, userID uint) (*entity.ParticipantPeriod, error)
}

type participantPeriodRepository struct {
	db *gorm.DB
}

func NewParticipantPeriodRepository(db *gorm.DB) ParticipantPeriodRepository {
	return &participantPeriodRepository{db: db}
}

func (r *participantPeriodRepository) Create(pp *entity.ParticipantPeriod) error {
	return r.db.Create(pp).Error
}

func (r *participantPeriodRepository) ExistsByPeriodAndUser(periodID, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&entity.ParticipantPeriod{}).
		Where("period_id = ? AND user_id = ?", periodID, userID).
		Count(&count).Error
	return count > 0, err
}

func (r *participantPeriodRepository) FindByPeriodAndUserPublicID(periodID uint, userPublicID string) (*entity.ParticipantPeriod, error) {
	parsedID, err := uuid.Parse(userPublicID)
	if err != nil {
		return nil, errors.New("invalid uuid format")
	}

	var pp entity.ParticipantPeriod
	err = r.db.
		Preload("User").
		Joins("JOIN users ON users.id = participant_periods.user_id").
		Where("participant_periods.period_id = ? AND users.public_id = ?", periodID, parsedID).
		First(&pp).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("participant not found in this period")
		}
		return nil, err
	}
	return &pp, nil
}

func (r *participantPeriodRepository) FindAllByPeriodIDPagination(periodID uint, filter, sort string, limit, offset int) ([]entity.ParticipantPeriod, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	var list []entity.ParticipantPeriod
	var total int64

	db := r.db.Model(&entity.ParticipantPeriod{}).
		Joins("JOIN users ON users.id = participant_periods.user_id").
		Where("participant_periods.period_id = ?", periodID)

	if filter != "" {
		filterPattern := "%" + filter + "%"
		db = db.Where("users.name ILIKE ? OR users.email ILIKE ? OR users.nim ILIKE ?", filterPattern, filterPattern, filterPattern)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	allowedSortFields := map[string]string{
		"status":      "participant_periods.status ASC",
		"-status":     "participant_periods.status DESC",
		"created_at":  "participant_periods.created_at ASC",
		"-created_at": "participant_periods.created_at DESC",
	}
	if sortClause, ok := allowedSortFields[sort]; ok {
		db = db.Order(sortClause)
	} else {
		db = db.Order("participant_periods.created_at DESC")
	}

	err := db.Preload("User").Limit(limit).Offset(offset).Find(&list).Error
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *participantPeriodRepository) Update(pp *entity.ParticipantPeriod) error {
	return r.db.Save(pp).Error
}

func (r *participantPeriodRepository) DeleteByPeriodAndUserPublicID(periodID uint, userPublicID string) error {
	parsedID, err := uuid.Parse(userPublicID)
	if err != nil {
		return errors.New("invalid uuid format")
	}

	result := r.db.
		Where("period_id = ? AND user_id = (SELECT id FROM users WHERE public_id = ?)", periodID, parsedID).
		Delete(&entity.ParticipantPeriod{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("participant not found in this period")
	}
	return nil
}

func (r *participantPeriodRepository) BulkCreate(list []entity.ParticipantPeriod) error {
	if len(list) == 0 {
		return nil
	}
	return r.db.Create(&list).Error
}

func (r *participantPeriodRepository) FindByUserID(userID uint) ([]entity.ParticipantPeriod, error) {
	var list []entity.ParticipantPeriod
	err := r.db.
		Preload("Period").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&list).Error
	return list, err
}

func (r *participantPeriodRepository) FindByPeriodIDAndUserID(periodID, userID uint) (*entity.ParticipantPeriod, error) {
	var pp entity.ParticipantPeriod
	err := r.db.
		Preload("Period").
		Preload("User").
		Where("period_id = ? AND user_id = ?", periodID, userID).
		First(&pp).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("participant not registered in this period")
		}
		return nil, err
	}
	return &pp, nil
}

func (r *participantPeriodRepository) FindByPeriodPublicIDAndUserID(periodPublicID string, userID uint) (*entity.ParticipantPeriod, error) {
	parsedPeriodID, err := uuid.Parse(periodPublicID)
	if err != nil {
		return nil, errors.New("invalid uuid format")
	}

	var pp entity.ParticipantPeriod
	err = r.db.
		Preload("Period").
		Preload("User").
		Joins("JOIN periods ON periods.id = participant_periods.period_id").
		Where("periods.public_id = ? AND participant_periods.user_id = ?", parsedPeriodID, userID).
		First(&pp).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("participant not registered in this period")
		}
		return nil, err
	}
	return &pp, nil
}
