package repository

import (
	"context"

	"wallet_api/internal/common/base"
	"wallet_api/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccountRepository struct {
	*base.BaseRepository[entity.Account] // Embed untuk Account (exposes base methods)
}

func New(db *gorm.DB) *AccountRepository {
	return &AccountRepository{
		BaseRepository: base.NewBaseRepository[entity.Account](db),
	}
}


func (r *AccountRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Account, error) {
	// Gunakan QueryBuilder dari base repo
	return r.NewQueryBuilder().
		Where("user_id", userID).
		Find(ctx)
}
