package exchangeRates

import (
	"context"

	"lemfi/simplebank/config"
	db "lemfi/simplebank/db/sqlc"
)

func (exchangeRateRepository *ExchangeRateRepository) ListExchangeRates(ctx context.Context) ([]db.ExchangeRate, error) {
	config.Logger.Info("Fetching all exchange rates")

	// Get all exchange rates from database
	dbExchangeRates, err := exchangeRateRepository.queries.ListExchangeRates(ctx)
	if err != nil {
		config.Logger.Error("Failed to fetch exchange rates", "error", err.Error())
		return []db.ExchangeRate{}, err
	}

	// Convert ListExchangeRatesRow to ExchangeRate
	exchangeRates := make([]db.ExchangeRate, len(dbExchangeRates))
	for i, rate := range dbExchangeRates {
		exchangeRates[i] = db.ExchangeRate{
			ID:           rate.ID,
			FromCurrency: rate.FromCurrency,
			ToCurrency:   rate.ToCurrency,
			Rate:         rate.Rate,
			CreatedAt:    rate.CreatedAt,
			UpdatedAt:    rate.UpdatedAt,
		}
	}

	config.Logger.Info("Successfully fetched exchange rates", "total", len(exchangeRates))
	return exchangeRates, nil
}
