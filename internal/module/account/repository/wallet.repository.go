package repository

import (
	"context"

	"wallet_api/internal/common/base"
	"wallet_api/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WalletRepository interface {
	base.Repository[entity.Wallet]
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Wallet, error)
}

type walletRepository struct {
	*base.BaseRepository[entity.Wallet]
}

func New(db *gorm.DB) WalletRepository {
	return &walletRepository{
		BaseRepository: base.NewBaseRepository[entity.Wallet](db),
	}
}

func (r *walletRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Wallet, error) {
	return r.NewQueryBuilder().
		Where("user_id", userID).
		Find(ctx)
}
