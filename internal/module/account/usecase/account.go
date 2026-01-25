package usecase

import (
	"context"
	"fmt"

	"wallet_api/internal/common/consts"
	"wallet_api/internal/common/errors"
	"wallet_api/internal/entity"

	"github.com/google/uuid"
)

type UseCase struct {
	accountRepo interface {
		Create(ctx context.Context, account *entity.Account) error
		FindByID(ctx context.Context, id uuid.UUID) (*entity.Account, error)
		FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Account, error)
		Update(ctx context.Context, account *entity.Account) error
	}
	transactionRepo interface {
		Create(ctx context.Context, transaction *entity.Transaction) error
		FindByAccountID(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]*entity.Transaction, error)
	}
}

func New(
	accountRepo interface {
		Create(ctx context.Context, account *entity.Account) error
		FindByID(ctx context.Context, id uuid.UUID) (*entity.Account, error)
		FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Account, error)
		Update(ctx context.Context, account *entity.Account) error
	},
	transactionRepo interface {
		Create(ctx context.Context, transaction *entity.Transaction) error
		FindByAccountID(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]*entity.Transaction, error)
	},
) *UseCase {
	return &UseCase{
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
	}
}

func (uc *UseCase) CreateAccount(ctx context.Context, userID uuid.UUID, accountName, currency string) (*entity.Account, error) {
	account := &entity.Account{
		UserID:      userID,
		AccountName: accountName,
		Balance:     0,
		Currency:    currency,
		Status:      consts.AccountStatusActive,
	}

	if err := uc.accountRepo.Create(ctx, account); err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	return account, nil
}

func (uc *UseCase) GetAccount(ctx context.Context, accountID uuid.UUID) (*entity.Account, error) {
	account, err := uc.accountRepo.FindByID(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}
	if account == nil {
		return nil, errors.ErrNotFound
	}

	return account, nil
}

func (uc *UseCase) GetUserAccounts(ctx context.Context, userID uuid.UUID) ([]*entity.Account, error) {
	accounts, err := uc.accountRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user accounts: %w", err)
	}

	return accounts, nil
}

func (uc *UseCase) Deposit(ctx context.Context, accountID uuid.UUID, amount int64, description string) error {
	if amount <= 0 {
		return errors.ErrBadRequest
	}

	// Get account
	account, err := uc.accountRepo.FindByID(ctx, accountID)
	if err != nil {
		return fmt.Errorf("failed to get account: %w", err)
	}
	if account == nil {
		return errors.ErrNotFound
	}

	// Update balance
	account.Balance += amount
	if err := uc.accountRepo.Update(ctx, account); err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}

	// Create transaction
	transaction := &entity.Transaction{
		AccountID:     accountID,
		ReferenceID:   uuid.New().String(),
		Type:          consts.TransactionTypeDeposit,
		Amount:        amount,
		BalanceBefore: account.Balance - amount,
		BalanceAfter:  account.Balance,
		Description:   description,
	}

	if err := uc.transactionRepo.Create(ctx, transaction); err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	return nil
}

func (uc *UseCase) Withdraw(ctx context.Context, accountID uuid.UUID, amount int64, description string) error {
	if amount <= 0 {
		return errors.ErrBadRequest
	}

	// Get account
	account, err := uc.accountRepo.FindByID(ctx, accountID)
	if err != nil {
		return fmt.Errorf("failed to get account: %w", err)
	}
	if account == nil {
		return errors.ErrNotFound
	}

	// Check balance
	if account.Balance < amount {
		return errors.New(400, "Insufficient balance", nil)
	}

	// Update balance
	account.Balance -= amount
	if err := uc.accountRepo.Update(ctx, account); err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}

	// Create transaction
	transaction := &entity.Transaction{
		AccountID:     accountID,
		ReferenceID:   uuid.New().String(),
		Type:          consts.TransactionTypeWithdrawal,
		Amount:        amount,
		BalanceBefore: account.Balance + amount,
		BalanceAfter:  account.Balance,
		Description:   description,
	}

	if err := uc.transactionRepo.Create(ctx, transaction); err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	return nil
}

func (uc *UseCase) GetTransactions(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]*entity.Transaction, error) {
	transactions, err := uc.transactionRepo.FindByAccountID(ctx, accountID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	return transactions, nil
}
