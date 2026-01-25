package entity

import (
	"time"

	"github.com/google/uuid"
)

// User represents user entity
type User struct {
	ID           uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Username     string         `json:"username" gorm:"uniqueIndex;not null;size:255"`
	PasswordHash string         `json:"-" gorm:"not null;size:255"`
	CreatedAt    time.Time      `json:"created_at"`
}

// TableName specifies the table name for User
func (User) TableName() string {
	return "users"
}
