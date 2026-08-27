package repositories

import (
	"errors"
	"siuji-backend/models"

	"gorm.io/gorm"
)

type ParticipantRepository interface {
	AssignParticipant(participantPeriod *models.ParticipantPeriod) error
	BulkAssignParticipants(participantPeriods []models.ParticipantPeriod) error
	FindByPeriodIDPagination(periodID uint, filter, sort string, limit, offset int) ([]models.ParticipantPeriod, int64, error)
	FindByPeriodAndUserPublicID(periodID uint, userPublicID string) (*models.ParticipantPeriod, error)
	Update(participantPeriod *models.ParticipantPeriod) error
	RemoveFromPeriod(periodID uint, userID uint) error
	FindByID(id uint) (*models.ParticipantPeriod, error)
	UpdateStatus(id uint, status string) error
	UpdateScore(id uint, score int) error
	// participant side
	FindByUserID(userID uint) ([]models.ParticipantPeriod, error)
	FindByPeriodIDAndUserID(periodID, userID uint) (*models.ParticipantPeriod, error)
}

type participantRepository struct {
	db *gorm.DB
}

func NewParticipantRepository(db *gorm.DB) ParticipantRepository {
	return &participantRepository{db: db}
}

func (r *participantRepository) AssignParticipant(participantPeriod *models.ParticipantPeriod) error {
	return r.db.Create(participantPeriod).Error
}

func (r *participantRepository) BulkAssignParticipants(participantPeriods []models.ParticipantPeriod) error {
	return r.db.CreateInBatches(&participantPeriods, 100).Error
}

func (r *participantRepository) FindByPeriodIDPagination(periodID uint, filter, sort string, limit, offset int) ([]models.ParticipantPeriod, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	var participants []models.ParticipantPeriod
	var total int64

	// Karena data peserta (User) ada di tabel relasi, kita gunakan Joins jika ingin memfilter berdasarkan data User (nama/email/nim)
	db := r.db.Model(&models.ParticipantPeriod{}).
		Joins("JOIN users ON users.id = participant_periods.user_id").
		Where("participant_periods.period_id = ?", periodID)

	// Jika ada parameter filter (misal mencari nama/email peserta di period tersebut)
	if filter != "" {
		filterPattern := "%" + filter + "%"
		db = db.Where("users.name ILIKE ? OR users.email ILIKE ? OR users.nim ILIKE ?", filterPattern, filterPattern, filterPattern)
	}

	// Hitung total data setelah filter
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Whitelist sorting field
	allowedSortFields := map[string]string{
		"id":          "participant_periods.id ASC",
		"-id":         "participant_periods.id DESC",
		"status":      "participant_periods.status ASC",
		"-status":     "participant_periods.status DESC",
		"score":       "participant_periods.score ASC",
		"-score":      "participant_periods.score DESC",
		"name":        "users.name ASC",
		"-name":       "users.name DESC",
		"created_at":  "participant_periods.created_at ASC",
		"-created_at": "participant_periods.created_at DESC",
	}

	if sort == "" {
		sort = "-created_at"
	}

	if sortClause, ok := allowedSortFields[sort]; ok {
		db = db.Order(sortClause)
	} else {
		db = db.Order("participant_periods.created_at DESC")
	}

	// Ambil data dengan Preload("User") agar objek user terisi lengkap di JSON response, serta batasi dengan Limit & Offset
	err := db.Preload("User").Limit(limit).Offset(offset).Find(&participants).Error
	if err != nil {
		return nil, 0, err
	}

	return participants, total, nil
}

func (r *participantRepository) FindByPeriodAndUserPublicID(periodID uint, userPublicID string) (*models.ParticipantPeriod, error) {
	var participantPeriod models.ParticipantPeriod
	err := r.db.Joins("User").
			Where("participant_periods.period_id = ? AND users.public_id = ?", periodID, userPublicID).
			First(&participantPeriod).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil,err
	}
	return &participantPeriod, nil
}

func (r *participantRepository) Update(participantPeriod *models.ParticipantPeriod) error {
	return r.db.Save(participantPeriod).Error
}

func (r *participantRepository) RemoveFromPeriod(periodID uint, userID uint) error {
	return r.db.Where("period_id = ? AND user_id = ?", periodID, userID).Delete(&models.ParticipantPeriod{}).Error
}

func (r *participantRepository) FindByID(id uint) (*models.ParticipantPeriod, error) {
	var participantPeriod models.ParticipantPeriod
	err := r.db.Preload("User").First(&participantPeriod, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &participantPeriod, nil
}

func (r *participantRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&models.ParticipantPeriod{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *participantRepository) UpdateScore(id uint, score int) error {
    return r.db.Model(&models.ParticipantPeriod{}).Where("id = ?", id).Update("score", score).Error
}

func (r *participantRepository) FindByUserID(userID uint) ([]models.ParticipantPeriod, error) {
	var participantPeriods []models.ParticipantPeriod
    err := r.db.Preload("Period").
        Where("user_id = ?", userID).
        Find(&participantPeriods).Error
    if err != nil {
        return nil, err
    }
    return participantPeriods, nil
}

func (r *participantRepository) FindByPeriodIDAndUserID(periodID, userID uint) (*models.ParticipantPeriod, error) {
	var participantPeriod models.ParticipantPeriod
    err := r.db.Preload("Period").
        Where("period_id = ? AND user_id = ?", periodID, userID).
        First(&participantPeriod).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &participantPeriod, nil
}