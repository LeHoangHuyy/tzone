package model

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;column:id"`
	UserID    uuid.UUID `gorm:"type:uuid;index;column:user_id"`
	ExpiresAt time.Time `gorm:"column:expires_at"`
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
