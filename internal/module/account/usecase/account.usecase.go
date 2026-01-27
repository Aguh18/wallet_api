package accountusecase

import (
	"context"
	"fmt"

	"wallet_api/internal/common/consts"
	"wallet_api/internal/common/errors"
	"wallet_api/internal/entity"
	"wallet_api/internal/module/account/repository"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type UseCase interface {
	CreateWallet(ctx context.Context, userID uuid.UUID, walletName, currency string) (*entity.Wallet, error)
	GetWallet(ctx context.Context, walletID uuid.UUID) (*entity.Wallet, error)
	GetUserWallets(ctx context.Context, userID uuid.UUID) ([]*entity.Wallet, error)
	Deposit(ctx context.Context, walletID uuid.UUID, amount decimal.Decimal, description string) error
	Withdraw(ctx context.Context, walletID uuid.UUID, amount decimal.Decimal, description string) error
	Transfer(ctx context.Context, fromWalletID, toWalletID uuid.UUID, amount decimal.Decimal, description string) error
	GetTransactions(ctx context.Context, walletID uuid.UUID, limit, offset int) ([]*entity.Transaction, error)
}

type useCase struct {
	walletRepo      repository.WalletRepository
	transactionRepo repository.TransactionRepository
}

func New(walletRepo repository.WalletRepository, transactionRepo repository.TransactionRepository) UseCase {
	return &useCase{
		walletRepo:      walletRepo,
		transactionRepo: transactionRepo,
	}
}

func (uc *useCase) CreateWallet(ctx context.Context, userID uuid.UUID, walletName, currency string) (*entity.Wallet, error) {
	wallet := &entity.Wallet{
		UserID:     userID,
		WalletName: walletName,
		Balance:    decimal.Zero,
		Currency:   currency,
		Status:     consts.WalletStatusActive,
	}

	if err := uc.walletRepo.Create(ctx, wallet); err != nil {
		return nil, fmt.Errorf("failed to create wallet: %w", err)
	}

	return wallet, nil
}

func (uc *useCase) GetWallet(ctx context.Context, walletID uuid.UUID) (*entity.Wallet, error) {
	wallet, err := uc.walletRepo.FindByID(ctx, walletID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}
	if wallet == nil {
		return nil, errors.ErrNotFound
	}

	return wallet, nil
}

func (uc *useCase) GetUserWallets(ctx context.Context, userID uuid.UUID) ([]*entity.Wallet, error) {
	wallets, err := uc.walletRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user wallets: %w", err)
	}

	return wallets, nil
}

func (uc *useCase) Deposit(ctx context.Context, walletID uuid.UUID, amount decimal.Decimal, description string) error {
	if amount.LessThanOrEqual(decimal.Zero) {
		return errors.ErrBadRequest
	}

	return uc.walletRepo.WithTransaction(ctx, func(tx *gorm.DB) error {
		// Get wallet with pessimistic locking
		wallet, err := uc.walletRepo.FindByIDForUpdate(ctx, walletID)
		if err != nil {
			return fmt.Errorf("failed to get wallet: %w", err)
		}
		if wallet == nil {
			return errors.ErrNotFound
		}

		// Calculate balance before and after
		balanceBefore := wallet.Balance
		balanceAfter := wallet.Balance.Add(amount)

		// Update balance
		wallet.Balance = balanceAfter
		if err := uc.walletRepo.Update(ctx, wallet); err != nil {
			return fmt.Errorf("failed to update wallet: %w", err)
		}

		// Create transaction
		transaction := &entity.Transaction{
			WalletID:      walletID,
			ReferenceID:   uuid.New().String(),
			Type:          consts.TransactionTypeDeposit,
			Amount:        amount,
			BalanceBefore: balanceBefore,
			BalanceAfter:  balanceAfter,
			Description:   description,
		}

		if err := uc.transactionRepo.Create(ctx, transaction); err != nil {
			return fmt.Errorf("failed to create transaction: %w", err)
		}

		return nil
	})
}

func (uc *useCase) Withdraw(ctx context.Context, walletID uuid.UUID, amount decimal.Decimal, description string) error {
	if amount.LessThanOrEqual(decimal.Zero) {
		return errors.ErrBadRequest
	}

	return uc.walletRepo.WithTransaction(ctx, func(tx *gorm.DB) error {
		// Get wallet with pessimistic locking
		wallet, err := uc.walletRepo.FindByIDForUpdate(ctx, walletID)
		if err != nil {
			return fmt.Errorf("failed to get wallet: %w", err)
		}
		if wallet == nil {
			return errors.ErrNotFound
		}

		// Check balance
		if wallet.Balance.LessThan(amount) {
			return errors.New(400, "Insufficient balance", nil)
		}

		// Calculate balance before and after
		balanceBefore := wallet.Balance
		balanceAfter := wallet.Balance.Sub(amount)

		// Update balance
		wallet.Balance = balanceAfter
		if err := uc.walletRepo.Update(ctx, wallet); err != nil {
			return fmt.Errorf("failed to update wallet: %w", err)
		}

		// Create transaction
		transaction := &entity.Transaction{
			WalletID:      walletID,
			ReferenceID:   uuid.New().String(),
			Type:          consts.TransactionTypeWithdrawal,
			Amount:        amount,
			BalanceBefore: balanceBefore,
			BalanceAfter:  balanceAfter,
			Description:   description,
		}

		if err := uc.transactionRepo.Create(ctx, transaction); err != nil {
			return fmt.Errorf("failed to create transaction: %w", err)
		}

		return nil
	})
}

func (uc *useCase) Transfer(ctx context.Context, fromWalletID, toWalletID uuid.UUID, amount decimal.Decimal, description string) error {
	if amount.LessThanOrEqual(decimal.Zero) {
		return errors.ErrBadRequest
	}

	if fromWalletID == toWalletID {
		return errors.New(400, "Cannot transfer to the same wallet", nil)
	}

	referenceID := uuid.New().String()

	return uc.walletRepo.WithTransaction(ctx, func(tx *gorm.DB) error {
		fromWallet, err := uc.walletRepo.FindByIDForUpdate(ctx, fromWalletID)
		if err != nil {
			return fmt.Errorf("failed to get from wallet: %w", err)
		}
		if fromWallet == nil {
			return errors.New(404, "Source wallet not found", nil)
		}

		toWallet, err := uc.walletRepo.FindByIDForUpdate(ctx, toWalletID)
		if err != nil {
			return fmt.Errorf("failed to get to wallet: %w", err)
		}
		if toWallet == nil {
			return errors.New(404, "Destination wallet not found", nil)
		}

		if fromWallet.Status != consts.WalletStatusActive {
			return errors.New(400, "Source wallet is not active", nil)
		}

		if toWallet.Status != consts.WalletStatusActive {
			return errors.New(400, "Destination wallet is not active", nil)
		}

		if fromWallet.Currency != toWallet.Currency {
			return errors.New(400, "Cannot transfer between different currencies", nil)
		}

		if fromWallet.Balance.LessThan(amount) {
			return errors.New(400, "Insufficient balance", nil)
		}

		// Calculate balances for from wallet
		fromBalanceBefore := fromWallet.Balance
		fromBalanceAfter := fromWallet.Balance.Sub(amount)
		fromWallet.Balance = fromBalanceAfter

		// Calculate balances for to wallet
		toBalanceBefore := toWallet.Balance
		toBalanceAfter := toWallet.Balance.Add(amount)
		toWallet.Balance = toBalanceAfter

		if err := uc.walletRepo.Update(ctx, fromWallet); err != nil {
			return fmt.Errorf("failed to update from wallet: %w", err)
		}

		if err := uc.walletRepo.Update(ctx, toWallet); err != nil {
			return fmt.Errorf("failed to update to wallet: %w", err)
		}

		withdrawalTx := &entity.Transaction{
			WalletID:      fromWalletID,
			ReferenceID:   referenceID,
			Type:          consts.TransactionTypeTransfer,
			Amount:        amount,
			BalanceBefore: fromBalanceBefore,
			BalanceAfter:  fromBalanceAfter,
			Description:   fmt.Sprintf("Transfer to wallet %s", toWalletID),
		}
		if description != "" {
			withdrawalTx.Description = fmt.Sprintf("%s - %s", description, withdrawalTx.Description)
		}

		if err := uc.transactionRepo.Create(ctx, withdrawalTx); err != nil {
			return fmt.Errorf("failed to create withdrawal transaction: %w", err)
		}

		depositTx := &entity.Transaction{
			WalletID:      toWalletID,
			ReferenceID:   referenceID,
			Type:          consts.TransactionTypeTransfer,
			Amount:        amount,
			BalanceBefore: toBalanceBefore,
			BalanceAfter:  toBalanceAfter,
			Description:   fmt.Sprintf("Transfer from wallet %s", fromWalletID),
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

func (uc *useCase) GetTransactions(ctx context.Context, walletID uuid.UUID, limit, offset int) ([]*entity.Transaction, error) {
	transactions, err := uc.transactionRepo.FindByWalletID(ctx, walletID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	return transactions, nil
}
