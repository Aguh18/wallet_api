package entity

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID      uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;index"`
	User        User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	AccountName string         `json:"account_name" gorm:"size:255"`
	Currency    string         `json:"currency" gorm:"default:'IDR';size:10"`
	Balance     int64          `json:"balance" gorm:"default:0;comment:Simpan dalam satuan terkecil, misal sen"`
	Status      string         `json:"status" gorm:"default:'active';size:50;comment:active, disabled"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func (Account) TableName() string {
	return "accounts"
}
