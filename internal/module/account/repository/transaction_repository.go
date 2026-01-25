package repository

import (
	"context"

	"wallet_api/internal/common/base"
	"wallet_api/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	*base.BaseRepository[entity.Transaction]
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{
		BaseRepository: base.NewBaseRepository[entity.Transaction](db),
	}
}


func (r *TransactionRepository) FindByAccountID(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]*entity.Transaction, error) {
	// Gunakan QueryBuilder
	return r.NewQueryBuilder().
		Where("account_id", accountID).
		Preload("Account").
		OrderBy("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(ctx)
}
