package exchangeRates

import (
	"context"

	"lemfi/simplebank/config"
	responses "lemfi/simplebank/internal/apps/exchangeRates/responses"
)

func (exchangeRateRepository *ExchangeRateRepository) ListExchangeRates(ctx context.Context) (responses.ListExchangeRatesResponse, error) {
	config.Logger.Info("Fetching all exchange rates")

	// Get all exchange rates from database
	dbExchangeRates, err := exchangeRateRepository.queries.ListExchangeRates(ctx)
	if err != nil {
		config.Logger.Error("Failed to fetch exchange rates", "error", err.Error())
		return responses.ListExchangeRatesResponse{}, err
	}

	// Convert database results to response format
	exchangeRates := make([]responses.ExchangeRateResponse, len(dbExchangeRates))
	for i, rate := range dbExchangeRates {
		exchangeRates[i] = responses.ExchangeRateResponse{
			ID:           rate.ID,
			FromCurrency: rate.FromCurrency,
			ToCurrency:   rate.ToCurrency,
			Rate:         rate.Rate,
			CreatedAt:    rate.CreatedAt.Time,
		}
	}

	response := responses.ListExchangeRatesResponse{
		ExchangeRates: exchangeRates,
		Total:         len(exchangeRates),
	}

	config.Logger.Info("Successfully fetched exchange rates", "total", response.Total)
	return response, nil
}
