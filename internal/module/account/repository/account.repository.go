package repository

import (
	"context"

	"wallet_api/internal/common/base"
	"wallet_api/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccountRepository interface {
	base.Repository[entity.Account]
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Account, error)
}

type accountRepository struct {
	*base.BaseRepository[entity.Account]
}

func New(db *gorm.DB) AccountRepository {
	return &accountRepository{
		BaseRepository: base.NewBaseRepository[entity.Account](db),
	}
}

func (r *accountRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Account, error) {
	return r.NewQueryBuilder().
		Where("user_id", userID).
		Find(ctx)
}
