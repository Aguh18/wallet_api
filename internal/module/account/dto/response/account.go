package response

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
