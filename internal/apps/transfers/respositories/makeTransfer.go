package transfers

import (
	"errors"

	"lemfi/simplebank/config"
	db "lemfi/simplebank/db/sqlc"
	transferErrors "lemfi/simplebank/internal/apps/transfers/errors"
	requests "lemfi/simplebank/internal/apps/transfers/requests"

	"github.com/shopspring/decimal"
)

func (transferRespository *TransferRespository) MakeTransfer(payload requests.MakeTransferRequest) (db.TransferTxResult, error) {
	// Validate that accounts exist and have sufficient balance
	fromAccount, err := transferRespository.queries.GetAccount(transferRespository.context, payload.FromAccountID)
	if err != nil {
		return db.TransferTxResult{}, transferErrors.ErrFromAccountNotFound
	}

	toAccount, err := transferRespository.queries.GetAccount(transferRespository.context, payload.ToAccountID)
	if err != nil {
		return db.TransferTxResult{}, transferErrors.ErrToAccountNotFound
	}

	// Validate currencies match (both fields are required)
	if fromAccount.Currency != payload.FromCurrency {
		config.Logger.Error("From account currency mismatch",
			"account_id", payload.FromAccountID,
			"account_currency", fromAccount.Currency,
			"requested_currency", payload.FromCurrency,
		)
		return db.TransferTxResult{}, transferErrors.ErrFromAccountCurrencyMismatch
	}

	if toAccount.Currency != payload.ToCurrency {
		config.Logger.Error("To account currency mismatch",
			"account_id", payload.ToAccountID,
			"account_currency", toAccount.Currency,
			"requested_currency", payload.ToCurrency,
		)
		return db.TransferTxResult{}, transferErrors.ErrToAccountCurrencyMismatch
	}

	// Validate sufficient balance
	if fromAccount.Balance.LessThan(payload.Amount) {
		return db.TransferTxResult{}, transferErrors.ErrInsufficientBalance
	}

	// Calculate converted amount and exchange rate
	convertedAmount := payload.Amount
	exchangeRate := decimal.NewFromInt(1) // Default to 1:1 for same currency

	// Calculate exchange rate for cross-currency transfers
	if payload.FromCurrency != payload.ToCurrency {
		// Get exchange rate from database
		rate, err := transferRespository.queries.GetExchangeRate(transferRespository.context, db.GetExchangeRateParams{
			FromCurrency: payload.FromCurrency,
			ToCurrency:   payload.ToCurrency,
		})
		if err != nil {
			config.Logger.Error("Failed to get exchange rate",
				"from_currency", payload.FromCurrency,
				"to_currency", payload.ToCurrency,
				"error", err.Error(),
			)
			return db.TransferTxResult{}, errors.New("exchange rate not found for currency pair")
		}

		exchangeRate = rate.Rate
		convertedAmount = payload.Amount.Mul(exchangeRate).Round(2)
	}

	// Prepare transfer parameters
	transferParams := db.TransferTxParams{
		FromAccountID:   payload.FromAccountID,
		ToAccountID:     payload.ToAccountID,
		Amount:          payload.Amount,
		ConvertedAmount: convertedAmount,
		ExchangeRate:    exchangeRate,
		FromCurrency:    payload.FromCurrency,
		ToCurrency:      payload.ToCurrency,
	}

	// Execute the transfer transaction
	result, err := transferRespository.queries.TransferTx(transferRespository.context, transferParams)
	if err != nil {
		return db.TransferTxResult{}, err
	}

	return result, nil
}
