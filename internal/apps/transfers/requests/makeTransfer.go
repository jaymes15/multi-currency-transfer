package transfers

type MakeTransferRequest struct {
	FromAccountID int64  `json:"from_account_id" validate:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" validate:"required,min=1"`
	Amount        int64  `json:"amount" validate:"required,min=1"`
	FromCurrency  string `json:"from_currency" validate:"required"`
	ToCurrency    string `json:"to_currency" validate:"required"`
}
