package response

import (
	"time"
	"wallet_api/internal/entity"
)

type WalletResponse struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	WalletName string `json:"wallet_name"`
	Currency   string `json:"currency"`
	Balance    string `json:"balance"`
	Status     string `json:"status"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type TransactionResponse struct {
	ID            string `json:"id"`
	WalletID      string `json:"wallet_id"`
	ReferenceID   string `json:"reference_id"`
	Type          string `json:"type"`
	Amount        string `json:"amount"`
	BalanceBefore string `json:"balance_before"`
	BalanceAfter  string `json:"balance_after"`
	Description   string `json:"description"`
	CreatedAt     string `json:"created_at"`
}

func ToWalletDto(wallet *entity.Wallet) WalletResponse {
	return WalletResponse{
		ID:         wallet.ID.String(),
		UserID:     wallet.UserID.String(),
		WalletName: wallet.WalletName,
		Currency:   wallet.Currency,
		Balance:    wallet.Balance.String(),
		Status:     wallet.Status,
		CreatedAt:  wallet.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  wallet.UpdatedAt.Format(time.RFC3339),
	}
}

func ToWalletDtos(wallets []*entity.Wallet) []WalletResponse {
	responses := make([]WalletResponse, len(wallets))
	for i, wallet := range wallets {
		responses[i] = ToWalletDto(wallet)
	}
	return responses
}

func ToTransactionDto(transaction *entity.Transaction) TransactionResponse {
	return TransactionResponse{
		ID:            transaction.ID.String(),
		WalletID:      transaction.WalletID.String(),
		ReferenceID:   transaction.ReferenceID,
		Type:          transaction.Type,
		Amount:        transaction.Amount.String(),
		BalanceBefore: transaction.BalanceBefore.String(),
		BalanceAfter:  transaction.BalanceAfter.String(),
		Description:   transaction.Description,
		CreatedAt:     transaction.CreatedAt.Format(time.RFC3339),
	}
}

func ToTransactionDtos(transactions []*entity.Transaction) []TransactionResponse {
	responses := make([]TransactionResponse, len(transactions))
	for i, transaction := range transactions {
		responses[i] = ToTransactionDto(transaction)
	}
	return responses
}
