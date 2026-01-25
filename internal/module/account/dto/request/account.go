package request

type CreateAccountRequest struct {
	AccountName string `json:"account_name" validate:"required"`
	Currency    string `json:"currency" validate:"required,default=IDR"`
}

type TransactionRequest struct {
	Amount      int64  `json:"amount" validate:"required,gt=0"`
	Description string `json:"description"`
}
