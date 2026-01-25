package repository

import (
	"context"
	"errors"

	"wallet_api/internal/common/base"
	"wallet_api/internal/entity"

	"gorm.io/gorm"
)

type UserRepository struct {
	*base.BaseRepository[entity.User] // Embed base (exposes all base methods!)
	db                                *gorm.DB
}

func New(db *gorm.DB) *UserRepository {
	return &UserRepository{
		BaseRepository: base.NewBaseRepository[entity.User](db),
		db:             db,
	}
}


func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
