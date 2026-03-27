package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"primaryKey;type:uuid;column:id"`
	Email        string    `gorm:"uniqueIndex;not null;column:email"`
	PasswordHash string    `gorm:"not null;column:password_hash"`
	CreatedAt    time.Time `gorm:"autoCreateTime;column:created_at"`
}

func (User) TableName() string {
	return "users"
}
