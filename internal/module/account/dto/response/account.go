package response

import (
	"time"
	"wallet_api/internal/entity"
)

type AccountResponse struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	AccountName string `json:"account_name"`
	Currency    string `json:"currency"`
	Balance     int64  `json:"balance"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type TransactionResponse struct {
	ID            string `json:"id"`
	AccountID     string `json:"account_id"`
	ReferenceID   string `json:"reference_id"`
	Type          string `json:"type"`
	Amount        int64  `json:"amount"`
	BalanceBefore int64  `json:"balance_before"`
	BalanceAfter  int64  `json:"balance_after"`
	Description   string `json:"description"`
	CreatedAt     string `json:"created_at"`
}

func ToAccountDto(account *entity.Account) AccountResponse {
	return AccountResponse{
		ID:          account.ID.String(),
		UserID:      account.UserID.String(),
		AccountName: account.AccountName,
		Currency:    account.Currency,
		Balance:     account.Balance,
		Status:      account.Status,
		CreatedAt:   account.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   account.UpdatedAt.Format(time.RFC3339),
	}
}

func ToAccountDtos(accounts []*entity.Account) []AccountResponse {
	responses := make([]AccountResponse, len(accounts))
	for i, account := range accounts {
		responses[i] = ToAccountDto(account)
	}
	return responses
}

func ToTransactionDto(transaction *entity.Transaction) TransactionResponse {
	return TransactionResponse{
		ID:            transaction.ID.String(),
		AccountID:     transaction.AccountID.String(),
		ReferenceID:   transaction.ReferenceID,
		Type:          transaction.Type,
		Amount:        transaction.Amount,
		BalanceBefore: transaction.BalanceBefore,
		BalanceAfter:  transaction.BalanceAfter,
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
