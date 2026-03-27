package repository

import (
	"github.com/LuuDinhTheTai/tzone/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type refreshTokenPostgresRepository struct {
	db *gorm.DB
}

// NewRefreshTokenRepository creates a new postgres-based refresh token repository.
func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return &refreshTokenPostgresRepository{db: db}
}

func (r *refreshTokenPostgresRepository) Create(token *model.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *refreshTokenPostgresRepository) FindByID(id uuid.UUID) (*model.RefreshToken, error) {
	var token model.RefreshToken
	err := r.db.Where("id = ?", id).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *refreshTokenPostgresRepository) DeleteByID(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&model.RefreshToken{}).Error
}

func (r *refreshTokenPostgresRepository) DeleteAllByUserID(userID uuid.UUID) error {
	return r.db.Where("user_id = ?", userID).Delete(&model.RefreshToken{}).Error
}
