package request

type CreateAccountRequest struct {
	AccountName string `json:"wallet_name" validate:"required"`
	Currency    string `json:"currency" validate:"required,default=IDR"`
}

type TransactionRequest struct {
	Amount      string `json:"amount" validate:"required,gt=0"`
	Description string `json:"description"`
}

type TransferRequest struct {
	ToWalletID string `json:"to_wallet_id" validate:"required"`
	Amount     string `json:"amount" validate:"required,gt=0"`
	Description string `json:"description"`
}
