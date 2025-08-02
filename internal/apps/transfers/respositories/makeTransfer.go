package transfers

import (
	"lemfi/simplebank/config"
	db "lemfi/simplebank/db/sqlc"
	transferErrors "lemfi/simplebank/internal/apps/transfers/errors"
	requests "lemfi/simplebank/internal/apps/transfers/requests"

	"github.com/shopspring/decimal"
)

func (transferRespository *TransferRespository) MakeTransfer(
	payload requests.MakeTransferRequest,
	convertedAmount decimal.Decimal,
	exchangeRate decimal.Decimal,
) (db.TransferTxResult, error) {
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

	// Prepare transfer parameters using values calculated in service layer
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
