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
	FindByWalletID(ctx context.Context, walletID uuid.UUID, limit, offset int) ([]*entity.Transaction, error)
}

type transactionRepository struct {
	*base.BaseRepository[entity.Transaction]
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{
		BaseRepository: base.NewBaseRepository[entity.Transaction](db),
	}
}

func (r *transactionRepository) FindByWalletID(ctx context.Context, walletID uuid.UUID, limit, offset int) ([]*entity.Transaction, error) {
	return r.NewQueryBuilder().
		Where("wallet_id", walletID).
		Preload("Wallet").
		OrderBy("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(ctx)
}
