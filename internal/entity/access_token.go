package entity

import (
	"time"

	"github.com/google/uuid"
)

// AccessToken represents access token entity
type AccessToken struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	SessionID uuid.UUID      `json:"session_id" gorm:"type:uuid;not null;index"`
	Session   Session        `json:"session,omitempty" gorm:"foreignKey:SessionID"`
	TokenHash string         `json:"-" gorm:"uniqueIndex;not null;size:500"`
	UserID    uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;index"`
	User      User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	ExpiredAt time.Time      `json:"expired_at" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
}

// TableName specifies the table name for AccessToken
func (AccessToken) TableName() string {
	return "access_tokens"
}
