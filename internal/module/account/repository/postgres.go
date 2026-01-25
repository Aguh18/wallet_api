package repository

import (
	"context"

	"github.com/google/uuid"
	"wallet_api/internal/entity"
	"gorm.io/gorm"
)

type postgresRepository struct {
	db *gorm.DB
}

// New creates new account repository
func New(db *gorm.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) CreateAccount(ctx context.Context, account *entity.Account) error {
	return r.db.WithContext(ctx).Create(account).Error
}

func (r *postgresRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Account, error) {
	var account entity.Account
	err := r.db.WithContext(ctx).Preload("User").First(&account, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &account, nil
}

func (r *postgresRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Account, error) {
	var accounts []entity.Account
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&accounts).Error
	return accounts, err
}

func (r *postgresRepository) UpdateAccount(ctx context.Context, account *entity.Account) error {
	return r.db.WithContext(ctx).Save(account).Error
}

func (r *postgresRepository) CreateTransaction(ctx context.Context, transaction *entity.Transaction) error {
	return r.db.WithContext(ctx).Create(transaction).Error
}

func (r *postgresRepository) FindByAccountID(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	err := r.db.WithContext(ctx).
		Preload("Account").
		Where("account_id = ?", accountID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error
	return transactions, err
}
