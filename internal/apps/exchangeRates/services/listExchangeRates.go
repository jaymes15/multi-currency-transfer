package exchangeRates

import (
	"context"

	"lemfi/simplebank/config"
	responses "lemfi/simplebank/internal/apps/exchangeRates/responses"
)

func (exchangeRateService *ExchangeRateService) ListExchangeRates(ctx context.Context) (responses.ListExchangeRatesResponse, error) {
	config.Logger.Info("Service: Listing all exchange rates")

	dbExchangeRates, err := exchangeRateService.exchangeRateRepository.ListExchangeRates(ctx)
	if err != nil {
		config.Logger.Error("Service: Failed to list exchange rates", "error", err.Error())
		return responses.ListExchangeRatesResponse{}, err
	}

	// Convert database results to response format
	exchangeRates := make([]responses.ExchangeRateResponse, len(dbExchangeRates))
	for i, rate := range dbExchangeRates {
		exchangeRates[i] = responses.NewExchangeRateResponse(rate)
	}

	config.Logger.Info("Service: Successfully listed exchange rates", "total", len(exchangeRates))
	return responses.ListExchangeRatesResponse{
		ExchangeRates: exchangeRates,
		Total:         len(exchangeRates),
	}, nil
}
