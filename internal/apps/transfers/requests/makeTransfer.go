package transfers

import "github.com/shopspring/decimal"

type MakeTransferRequest struct {
	FromAccountID int64           `json:"from_account_id" validate:"required,min=1"`
	ToAccountID   int64           `json:"to_account_id" validate:"required,min=1"`
	Amount        decimal.Decimal `json:"amount" validate:"required"`
	FromCurrency  string          `json:"from_currency" validate:"required"`
	ToCurrency    string          `json:"to_currency" validate:"required"`
	ExchangeRate  decimal.Decimal `json:"exchange_rate" validate:"omitempty"`
}
