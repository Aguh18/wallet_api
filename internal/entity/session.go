package entity

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID           uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID       uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;index"`
	User         User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	SessionToken string         `json:"session_token" gorm:"uniqueIndex;not null;size:500"`
	UserAgent    string         `json:"user_agent" gorm:"type:text"`
	IPAddress    string         `json:"ip_address" gorm:"size:45"`
	IsRevoked    bool           `json:"is_revoked" gorm:"default:false"`
	ExpiredAt    time.Time      `json:"expired_at" gorm:"not null;index"`
	CreatedAt    time.Time      `json:"created_at"`
}

func (Session) TableName() string {
	return "sessions"
}
