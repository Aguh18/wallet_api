package repository

import (
	"context"

	"wallet_api/internal/common/base"
	"wallet_api/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(ctx context.Context, transaction *entity.Transaction) error
	FindByAccountID(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]*entity.Transaction, error)
}

type transactionRepository struct {
	*base.BaseRepository[entity.Transaction]
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{
		BaseRepository: base.NewBaseRepository[entity.Transaction](db),
	}
}

func (r *transactionRepository) FindByAccountID(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]*entity.Transaction, error) {
	return r.NewQueryBuilder().
		Where("account_id", accountID).
		Preload("Account").
		OrderBy("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(ctx)
}
