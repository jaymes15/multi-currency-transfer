package responses

import (
	"time"

	db "lemfi/simplebank/db/sqlc"

	"github.com/shopspring/decimal"
)

type MakeTransferResponse struct {
	Transfer    TransferDetail `json:"transfer"`
	FromAccount AccountDetail  `json:"from_account"`
	ToAccount   AccountDetail  `json:"to_account"`
	FromEntry   EntryDetail    `json:"from_entry"`
	ToEntry     EntryDetail    `json:"to_entry"`
	Message     string         `json:"message"`
}

type TransferDetail struct {
	ID              int64           `json:"id"`
	FromAccountID   int64           `json:"from_account_id"`
	ToAccountID     int64           `json:"to_account_id"`
	Amount          decimal.Decimal `json:"amount"`
	ConvertedAmount decimal.Decimal `json:"converted_amount,omitempty"`
	FromCurrency    string          `json:"from_currency,omitempty"`
	ToCurrency      string          `json:"to_currency,omitempty"`
	ExchangeRate    decimal.Decimal `json:"exchange_rate,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
}

type AccountDetail struct {
	ID        int64           `json:"id"`
	Owner     string          `json:"owner"`
	Balance   decimal.Decimal `json:"balance"`
	Currency  string          `json:"currency"`
	CreatedAt time.Time       `json:"created_at"`
}

type EntryDetail struct {
	ID        int64           `json:"id"`
	AccountID int64           `json:"account_id"`
	Amount    decimal.Decimal `json:"amount"`
	CreatedAt time.Time       `json:"created_at"`
}

// NewMakeTransferResponse converts database result to API response
func NewMakeTransferResponse(result db.TransferTxResult) MakeTransferResponse {
	response := MakeTransferResponse{
		Transfer: TransferDetail{
			ID:              result.Transfer.ID,
			FromAccountID:   result.Transfer.FromAccountID,
			ToAccountID:     result.Transfer.ToAccountID,
			Amount:          result.Transfer.Amount,
			ConvertedAmount: result.Transfer.ConvertedAmount,
			ExchangeRate:    result.Transfer.ExchangeRate,
			CreatedAt:       result.Transfer.CreatedAt,
		},
		FromAccount: AccountDetail{
			ID:        result.FromAccount.ID,
			Owner:     result.FromAccount.Owner,
			Balance:   result.FromAccount.Balance,
			Currency:  result.FromAccount.Currency,
			CreatedAt: result.FromAccount.CreatedAt,
		},
		ToAccount: AccountDetail{
			ID:        result.ToAccount.ID,
			Owner:     result.ToAccount.Owner,
			Balance:   result.ToAccount.Balance,
			Currency:  result.ToAccount.Currency,
			CreatedAt: result.ToAccount.CreatedAt,
		},
		FromEntry: EntryDetail{
			ID:        result.FromEntry.ID,
			AccountID: result.FromEntry.AccountID,
			Amount:    result.FromEntry.Amount,
			CreatedAt: result.FromEntry.CreatedAt,
		},
		ToEntry: EntryDetail{
			ID:        result.ToEntry.ID,
			AccountID: result.ToEntry.AccountID,
			Amount:    result.ToEntry.Amount,
			CreatedAt: result.ToEntry.CreatedAt,
		},
		Message: "Transfer completed successfully",
	}

	// Add currency information if available
	if result.Transfer.FromCurrency.Valid {
		response.Transfer.FromCurrency = result.Transfer.FromCurrency.String
	}
	if result.Transfer.ToCurrency.Valid {
		response.Transfer.ToCurrency = result.Transfer.ToCurrency.String
	}

	return response
}
