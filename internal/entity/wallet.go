package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Wallet struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID      uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;index"`
	User        User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	WalletName  string         `json:"wallet_name" gorm:"size:255"`
	Currency    string         `json:"currency" gorm:"default:'IDR';size:10"`
	Balance     decimal.Decimal `json:"balance" gorm:"type:numeric(20,2);default:0"`
	Status      string         `json:"status" gorm:"default:'active';size:50;comment:active, disabled"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func (Wallet) TableName() string {
	return "wallets"
}
