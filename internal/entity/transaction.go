package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Transaction struct {
	ID            uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	WalletID      uuid.UUID      `json:"wallet_id" gorm:"type:uuid;not null;index"`
	Wallet        Wallet         `json:"wallet,omitempty" gorm:"foreignKey:WalletID"`
	ReferenceID   string         `json:"reference_id" gorm:"uniqueIndex;not null;size:500;comment:Untuk idempotency key"`
	Type          string         `json:"type" gorm:"not null;size:50;comment:deposit, withdrawal, transfer, payment"`
	Amount        decimal.Decimal `json:"amount" gorm:"type:numeric(20,2);not null"`
	BalanceBefore decimal.Decimal `json:"balance_before" gorm:"type:numeric(20,2);not null"`
	BalanceAfter  decimal.Decimal `json:"balance_after" gorm:"type:numeric(20,2);not null"`
	Description   string         `json:"description" gorm:"type:text"`
	CreatedAt     time.Time      `json:"created_at" gorm:"index"`
}

func (Transaction) TableName() string {
	return "transactions"
}
