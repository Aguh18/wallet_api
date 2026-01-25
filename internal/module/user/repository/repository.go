package repository

import (
	"context"

	"wallet_api/internal/entity"
	"github.com/google/uuid"
)

// Repository defines user data operations
type Repository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByUsername(ctx context.Context, username string) (*entity.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	List(ctx context.Context, limit, offset int) ([]entity.User, error)
}
