package repository

import (
	"github.com/LuuDinhTheTai/tzone/internal/model"
	"github.com/google/uuid"
)

// RefreshTokenRepository defines the necessary database operations for tracking refresh tokens.
type RefreshTokenRepository interface {
	Create(token *model.RefreshToken) error
	FindByID(id uuid.UUID) (*model.RefreshToken, error)
	DeleteByID(id uuid.UUID) error
	DeleteAllByUserID(userID uuid.UUID) error
}
