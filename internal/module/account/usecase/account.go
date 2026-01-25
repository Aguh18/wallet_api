package usecase

import (
	"context"
	"fmt"

	"wallet_api/internal/common/errors"
	"wallet_api/internal/common/consts"
	"wallet_api/internal/entity"
	"wallet_api/internal/module/account/repository"
	"github.com/google/uuid"
)

type UseCase struct {
	repo repository.Repository
}

// New creates new account usecase
func New(repo repository.Repository) *UseCase {
	return &UseCase{
		repo: repo,
	}
}

// CreateAccount creates new account for user
func (uc *UseCase) CreateAccount(ctx context.Context, userID uuid.UUID, accountName, currency string) (*entity.Account, error) {
	account := &entity.Account{
		UserID:      userID,
		AccountName: accountName,
		Balance:     0,
		Currency:    currency,
		Status:      consts.AccountStatusActive,
	}

	if err := uc.repo.CreateAccount(ctx, account); err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	return account, nil
}

// GetAccount retrieves account by ID
func (uc *UseCase) GetAccount(ctx context.Context, accountID uuid.UUID) (*entity.Account, error) {
	account, err := uc.repo.FindByID(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}
	if account == nil {
		return nil, errors.ErrNotFound
	}

	return account, nil
}

// GetUserAccounts retrieves all accounts for user
func (uc *UseCase) GetUserAccounts(ctx context.Context, userID uuid.UUID) ([]entity.Account, error) {
	accounts, err := uc.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user accounts: %w", err)
	}

	return accounts, nil
}

// Deposit adds money to account
func (uc *UseCase) Deposit(ctx context.Context, accountID uuid.UUID, amount int64, description string) error {
	if amount <= 0 {
		return errors.ErrBadRequest
	}

	// Get account
	account, err := uc.repo.FindByID(ctx, accountID)
	if err != nil {
		return fmt.Errorf("failed to get account: %w", err)
	}
	if account == nil {
		return errors.ErrNotFound
	}

	// Update balance
	account.Balance += amount
	if err := uc.repo.UpdateAccount(ctx, account); err != nil {
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

	if err := uc.repo.CreateTransaction(ctx, transaction); err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	return nil
}

// Withdraw removes money from account
func (uc *UseCase) Withdraw(ctx context.Context, accountID uuid.UUID, amount int64, description string) error {
	if amount <= 0 {
		return errors.ErrBadRequest
	}

	// Get account
	account, err := uc.repo.FindByID(ctx, accountID)
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
	if err := uc.repo.UpdateAccount(ctx, account); err != nil {
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

	if err := uc.repo.CreateTransaction(ctx, transaction); err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	return nil
}

// GetTransactions retrieves transaction history for account
func (uc *UseCase) GetTransactions(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]entity.Transaction, error) {
	transactions, err := uc.repo.FindByAccountID(ctx, accountID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	return transactions, nil
}
