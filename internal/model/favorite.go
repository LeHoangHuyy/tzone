package model

import (
	"time"

	"github.com/google/uuid"
)

type Favorite struct {
	ID        uint      `gorm:"primaryKey;autoIncrement;column:id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_favorites_user_device;column:user_id"`
	DeviceID  string    `gorm:"type:varchar(64);not null;uniqueIndex:idx_favorites_user_device;column:device_id"`
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at"`
}

func (Favorite) TableName() string {
	return "favorites"
}
