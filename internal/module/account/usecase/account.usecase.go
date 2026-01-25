package accountusecase

import (
	"context"
	"fmt"

	"wallet_api/internal/common/consts"
	"wallet_api/internal/common/errors"
	"wallet_api/internal/entity"
	"wallet_api/internal/module/account/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UseCase interface {
	CreateAccount(ctx context.Context, userID uuid.UUID, accountName, currency string) (*entity.Account, error)
	GetAccount(ctx context.Context, accountID uuid.UUID) (*entity.Account, error)
	GetUserAccounts(ctx context.Context, userID uuid.UUID) ([]*entity.Account, error)
	Deposit(ctx context.Context, accountID uuid.UUID, amount int64, description string) error
	Withdraw(ctx context.Context, accountID uuid.UUID, amount int64, description string) error
	Transfer(ctx context.Context, fromAccountID, toAccountID uuid.UUID, amount int64, description string) error
	GetTransactions(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]*entity.Transaction, error)
}

type useCase struct {
	accountRepo     repository.AccountRepository
	transactionRepo repository.TransactionRepository
}

func New(accountRepo repository.AccountRepository, transactionRepo repository.TransactionRepository) UseCase {
	return &useCase{
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
	}
}

func (uc *useCase) CreateAccount(ctx context.Context, userID uuid.UUID, accountName, currency string) (*entity.Account, error) {
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

func (uc *useCase) GetAccount(ctx context.Context, accountID uuid.UUID) (*entity.Account, error) {
	account, err := uc.accountRepo.FindByID(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}
	if account == nil {
		return nil, errors.ErrNotFound
	}

	return account, nil
}

func (uc *useCase) GetUserAccounts(ctx context.Context, userID uuid.UUID) ([]*entity.Account, error) {
	accounts, err := uc.accountRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user accounts: %w", err)
	}

	return accounts, nil
}

func (uc *useCase) Deposit(ctx context.Context, accountID uuid.UUID, amount int64, description string) error {
	if amount <= 0 {
		return errors.ErrBadRequest
	}

	return uc.accountRepo.WithTransaction(ctx, func(tx *gorm.DB) error {
		// Get account with pessimistic locking
		account, err := uc.accountRepo.FindByIDForUpdate(ctx, accountID)
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
	})
}

func (uc *useCase) Withdraw(ctx context.Context, accountID uuid.UUID, amount int64, description string) error {
	if amount <= 0 {
		return errors.ErrBadRequest
	}

	return uc.accountRepo.WithTransaction(ctx, func(tx *gorm.DB) error {
		// Get account with pessimistic locking
		account, err := uc.accountRepo.FindByIDForUpdate(ctx, accountID)
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
	})
}

func (uc *useCase) Transfer(ctx context.Context, fromAccountID, toAccountID uuid.UUID, amount int64, description string) error {
	if amount <= 0 {
		return errors.ErrBadRequest
	}

	if fromAccountID == toAccountID {
		return errors.New(400, "Cannot transfer to the same account", nil)
	}

	referenceID := uuid.New().String()

	return uc.accountRepo.WithTransaction(ctx, func(tx *gorm.DB) error {
		fromAccount, err := uc.accountRepo.FindByIDForUpdate(ctx, fromAccountID)
		if err != nil {
			return fmt.Errorf("failed to get from account: %w", err)
		}
		if fromAccount == nil {
			return errors.New(404, "Source account not found", nil)
		}

		toAccount, err := uc.accountRepo.FindByIDForUpdate(ctx, toAccountID)
		if err != nil {
			return fmt.Errorf("failed to get to account: %w", err)
		}
		if toAccount == nil {
			return errors.New(404, "Destination account not found", nil)
		}

		if fromAccount.Status != consts.AccountStatusActive {
			return errors.New(400, "Source account is not active", nil)
		}

		if toAccount.Status != consts.AccountStatusActive {
			return errors.New(400, "Destination account is not active", nil)
		}

		if fromAccount.Currency != toAccount.Currency {
			return errors.New(400, "Cannot transfer between different currencies", nil)
		}

		if fromAccount.Balance < amount {
			return errors.New(400, "Insufficient balance", nil)
		}

		fromAccount.Balance -= amount
		if err := uc.accountRepo.Update(ctx, fromAccount); err != nil {
			return fmt.Errorf("failed to update from account: %w", err)
		}

		toAccount.Balance += amount
		if err := uc.accountRepo.Update(ctx, toAccount); err != nil {
			return fmt.Errorf("failed to update to account: %w", err)
		}

		withdrawalTx := &entity.Transaction{
			AccountID:     fromAccountID,
			ReferenceID:   referenceID,
			Type:          consts.TransactionTypeTransfer,
			Amount:        amount,
			BalanceBefore: fromAccount.Balance + amount,
			BalanceAfter:  fromAccount.Balance,
			Description:   fmt.Sprintf("Transfer to account %s", toAccountID),
		}
		if description != "" {
			withdrawalTx.Description = fmt.Sprintf("%s - %s", description, withdrawalTx.Description)
		}

		if err := uc.transactionRepo.Create(ctx, withdrawalTx); err != nil {
			return fmt.Errorf("failed to create withdrawal transaction: %w", err)
		}

		depositTx := &entity.Transaction{
			AccountID:     toAccountID,
			ReferenceID:   referenceID,
			Type:          consts.TransactionTypeTransfer,
			Amount:        amount,
			BalanceBefore: toAccount.Balance - amount,
			BalanceAfter:  toAccount.Balance,
			Description:   fmt.Sprintf("Transfer from account %s", fromAccountID),
		}
		if description != "" {
			depositTx.Description = fmt.Sprintf("%s - %s", description, depositTx.Description)
		}

		if err := uc.transactionRepo.Create(ctx, depositTx); err != nil {
			return fmt.Errorf("failed to create deposit transaction: %w", err)
		}

		return nil
	})
}

func (uc *useCase) GetTransactions(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]*entity.Transaction, error) {
	transactions, err := uc.transactionRepo.FindByAccountID(ctx, accountID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	return transactions, nil
}
