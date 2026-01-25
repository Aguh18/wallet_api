package repository

import (
	"context"

	"wallet_api/internal/entity"
	"github.com/google/uuid"
)

// Repository defines account data operations
type Repository interface {
	// Account operations
	CreateAccount(ctx context.Context, account *entity.Account) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Account, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Account, error)
	UpdateAccount(ctx context.Context, account *entity.Account) error

	// Transaction operations
	CreateTransaction(ctx context.Context, transaction *entity.Transaction) error
	FindByAccountID(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]entity.Transaction, error)
}
