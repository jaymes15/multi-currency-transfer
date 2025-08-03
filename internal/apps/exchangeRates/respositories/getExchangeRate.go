package exchangeRates

import (
	"context"

	"lemfi/simplebank/config"
	db "lemfi/simplebank/db/sqlc"
	exchangeRateErrors "lemfi/simplebank/internal/apps/exchangeRates/errors"
	requests "lemfi/simplebank/internal/apps/exchangeRates/requests"
)

func (exchangeRateRepository *ExchangeRateRepository) GetExchangeRate(ctx context.Context, payload requests.GetExchangeRateRequest) (db.ExchangeRate, error) {
	config.Logger.Info("Fetching exchange rate for currency pair",
		"from_currency", payload.FromCurrency,
		"to_currency", payload.ToCurrency,
		"amount", payload.Amount.String(),
	)

	// Get exchange rate from database
	exchangeRate, err := exchangeRateRepository.queries.GetExchangeRate(ctx, db.GetExchangeRateParams{
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

	// Convert GetExchangeRateRow to ExchangeRate
	return db.ExchangeRate{
		ID:           exchangeRate.ID,
		FromCurrency: exchangeRate.FromCurrency,
		ToCurrency:   exchangeRate.ToCurrency,
		Rate:         exchangeRate.Rate,
		CreatedAt:    exchangeRate.CreatedAt,
		UpdatedAt:    exchangeRate.UpdatedAt,
	}, nil
}
