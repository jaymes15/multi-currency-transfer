package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountID   int64           `json:"from_account_id"`
	ToAccountID     int64           `json:"to_account_id"`
	Amount          decimal.Decimal `json:"amount"`
	ConvertedAmount decimal.Decimal `json:"converted_amount"`
	ExchangeRate    decimal.Decimal `json:"exchange_rate,omitempty"`
	FromCurrency    string          `json:"from_currency,omitempty"`
	ToCurrency      string          `json:"to_currency,omitempty"`
	Fee             decimal.Decimal `json:"fee,omitempty"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx performs a money transfer from one account to the other.
// It creates the transfer, add account entries, and update accounts' balance within a database transaction
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// Validate that accounts exist
		fromAccount, err := q.GetAccount(ctx, arg.FromAccountID)
		if err != nil {
			return fmt.Errorf("from account not found: %w", err)
		}

		toAccount, err := q.GetAccount(ctx, arg.ToAccountID)
		if err != nil {
			return fmt.Errorf("to account not found: %w", err)
		}

		// Validate currencies match if provided
		if arg.FromCurrency != "" && fromAccount.Currency != arg.FromCurrency {
			return fmt.Errorf("from account currency mismatch: expected %s, got %s", fromAccount.Currency, arg.FromCurrency)
		}

		if arg.ToCurrency != "" && toAccount.Currency != arg.ToCurrency {
			return fmt.Errorf("to account currency mismatch: expected %s, got %s", toAccount.Currency, arg.ToCurrency)
		}

		// Validate sufficient balance (including fee)
		totalAmount := arg.Amount.Add(arg.Fee)
		if fromAccount.Balance.LessThan(totalAmount) {
			return fmt.Errorf("insufficient balance: account has %s, transfer requires %s (amount: %s + fee: %s)", fromAccount.Balance.String(), totalAmount.String(), arg.Amount.String(), arg.Fee.String())
		}

		// Prepare currency data
		var fromCurrency pgtype.Text
		var toCurrency pgtype.Text

		if arg.FromCurrency != "" {
			fromCurrency.Scan(arg.FromCurrency)
		}
		if arg.ToCurrency != "" {
			toCurrency.Scan(arg.ToCurrency)
		}

		createTransferResult, err := q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID:   arg.FromAccountID,
			ToAccountID:     arg.ToAccountID,
			Amount:          arg.Amount,
			ConvertedAmount: arg.ConvertedAmount,
			ExchangeRate:    arg.ExchangeRate,
			FromCurrency:    fromCurrency,
			ToCurrency:      toCurrency,
			Fee:             arg.Fee,
		})
		if err != nil {
			return err
		}

		// Convert CreateTransferRow to Transfer
		result.Transfer = Transfer{
			ID:              createTransferResult.ID,
			FromAccountID:   createTransferResult.FromAccountID,
			ToAccountID:     createTransferResult.ToAccountID,
			Amount:          createTransferResult.Amount,
			ConvertedAmount: createTransferResult.ConvertedAmount,
			ExchangeRate:    createTransferResult.ExchangeRate,
			FromCurrency:    createTransferResult.FromCurrency,
			ToCurrency:      createTransferResult.ToCurrency,
			Fee:             createTransferResult.Fee,
			CreatedAt:       createTransferResult.CreatedAt,
		}

		// Create entries with original amount + fee (from account) and converted amount (to account)
		fromEntryAmount := arg.Amount.Add(arg.Fee).Neg() // Debit original amount + fee (negative)
		toEntryAmount := arg.ConvertedAmount             // Credit converted amount

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    fromEntryAmount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    toEntryAmount,
		})
		if err != nil {
			return err
		}

		// Update account balances with appropriate amounts
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, fromEntryAmount, arg.ToAccountID, toEntryAmount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, toEntryAmount, arg.FromAccountID, fromEntryAmount)
		}

		return err
	})

	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 decimal.Decimal,
	accountID2 int64,
	amount2 decimal.Decimal,
) (account1 Account, account2 Account, err error) {
	_, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	_, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	if err != nil {
		return
	}

	// Get the updated accounts
	account1, err = q.GetAccount(ctx, accountID1)
	if err != nil {
		return
	}

	account2, err = q.GetAccount(ctx, accountID2)
	return
}
