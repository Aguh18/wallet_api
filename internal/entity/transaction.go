package entity

import (
	"time"

	"github.com/google/uuid"
)

// Transaction represents transaction entity
type Transaction struct {
	ID            uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	AccountID     uuid.UUID      `json:"account_id" gorm:"type:uuid;not null;index"`
	Account       Account        `json:"account,omitempty" gorm:"foreignKey:AccountID"`
	ReferenceID   string         `json:"reference_id" gorm:"uniqueIndex;not null;size:500;comment:Untuk idempotency key"`
	Type          string         `json:"type" gorm:"not null;size:50;comment:deposit, withdrawal, transfer, payment"`
	Amount        int64          `json:"amount" gorm:"not null"`
	BalanceBefore int64          `json:"balance_before" gorm:"not null"`
	BalanceAfter  int64          `json:"balance_after" gorm:"not null"`
	Description   string         `json:"description" gorm:"type:text"`
	CreatedAt     time.Time      `json:"created_at" gorm:"index"`
}

// TableName specifies the table name for Transaction
func (Transaction) TableName() string {
	return "transactions"
}
