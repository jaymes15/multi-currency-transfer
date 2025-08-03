package exchangeRates

import (
	"context"

	"lemfi/simplebank/config"
	db "lemfi/simplebank/db/sqlc"
	"lemfi/simplebank/internal/apps/currencies"
	exchangeRateErrors "lemfi/simplebank/internal/apps/exchangeRates/errors"
	requests "lemfi/simplebank/internal/apps/exchangeRates/requests"

	"github.com/shopspring/decimal"
)

func (exchangeRateRepository *ExchangeRateRepository) GetExchangeRate(ctx context.Context, payload requests.GetExchangeRateRequest) (db.ExchangeRate, error) {
	config.Logger.Info("Fetching exchange rate for currency pair",
		"from_currency", payload.FromCurrency,
		"to_currency", payload.ToCurrency,
		"amount", payload.Amount.String(),
	)

	// Validate currencies are supported
	if !currencies.IsSupportedCurrency(currencies.Currency(payload.FromCurrency)) {
		config.Logger.Error("From currency is not supported", "currency", payload.FromCurrency)
		return db.ExchangeRate{}, exchangeRateErrors.ErrUnsupportedCurrency
	}

	if !currencies.IsSupportedCurrency(currencies.Currency(payload.ToCurrency)) {
		config.Logger.Error("To currency is not supported", "currency", payload.ToCurrency)
		return db.ExchangeRate{}, exchangeRateErrors.ErrUnsupportedCurrency
	}

	// Validate amount is positive
	if payload.Amount.LessThanOrEqual(decimal.Zero) {
		config.Logger.Error("Invalid amount", "amount", payload.Amount.String())
		return db.ExchangeRate{}, exchangeRateErrors.ErrInvalidAmount
	}

	// Get exchange rate from database
	dbExchangeRate, err := exchangeRateRepository.queries.GetExchangeRate(ctx, db.GetExchangeRateParams{
		FromCurrency: payload.FromCurrency,
		ToCurrency:   payload.ToCurrency,
	})
	if err != nil {
		config.Logger.Error("Failed to get exchange rate",
			"from_currency", payload.FromCurrency,
			"to_currency", payload.ToCurrency,
			"error", err.Error(),
		)
		return db.ExchangeRate{}, exchangeRateErrors.ErrExchangeRateNotFound
	}

	return dbExchangeRate, nil
}
