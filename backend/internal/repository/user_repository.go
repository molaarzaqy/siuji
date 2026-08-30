package repository

import (
	"errors"
	"siuji-backend/internal/entity"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *entity.User) error
	FindByID(id uint) (*entity.User, error)
	FindByPublicID(publicID string) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
	UpdatePassword(userID uint, hashedPassword string) error
	FindAllPagination(filter, sort string, limit, offset int) ([]entity.User, int64, error)
	Delete(publicID string) error
	Update(user *entity.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *entity.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByID(id uint) (*entity.User, error) {
	var user entity.User
	err := r.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByPublicID(publicID string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("public_id = ?", publicID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdatePassword(userID uint, hashedPassword string) error {
	return r.db.Model(&entity.User{}).
		Where("id = ?", userID).Update("password", hashedPassword).Error
}

func (r *userRepository) FindAllPagination(filter, sort string, limit, offset int) ([]entity.User, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	var users []entity.User
	var total int64

	db := r.db.Model(&entity.User{})
	if filter != "" {
		filterPattern := "%" + filter + "%"
		db = db.Where("name ILIKE ? OR email ILIKE ? OR nim ILIKE ? OR university ILIKE ?",
			filterPattern, filterPattern, filterPattern, filterPattern)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	allowedSortFields := map[string]string{
		"id":          "id ASC",
		"-id":         "id DESC",
		"name":        "name ASC",
		"-name":       "name DESC",
		"email":       "email ASC",
		"-email":      "email DESC",
		"nim":         "nim ASC",
		"-nim":        "nim DESC",
		"role":        "role ASC",
		"-role":       "role DESC",
		"created_at":  "created_at ASC",
		"-created_at": "created_at DESC",
	}

	if sort == "" {
		sort = "-created_at"
	}
	if sortClause, ok := allowedSortFields[sort]; ok {
		db = db.Order(sortClause)
	} else {
		db = db.Order("created_at DESC")
	}

	err := db.Limit(limit).Offset(offset).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *userRepository) Delete(publicID string) error {
	result := r.db.Where("public_id = ?", publicID).Delete(&entity.User{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *userRepository) Update(user *entity.User) error {
	return r.db.Save(user).Error
}
