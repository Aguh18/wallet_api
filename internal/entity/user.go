package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Username     string         `json:"username" gorm:"uniqueIndex;not null;size:255"`
	Email        string         `json:"email" gorm:"uniqueIndex;size:255"`
	PasswordHash string         `json:"-" gorm:"not null;size:255"`
	CreatedAt    time.Time      `json:"created_at"`
}

func (User) TableName() string {
	return "users"
}
