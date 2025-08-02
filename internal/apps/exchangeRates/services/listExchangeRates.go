package exchangeRates

import (
	"context"

	"lemfi/simplebank/config"
	responses "lemfi/simplebank/internal/apps/exchangeRates/responses"
)

func (exchangeRateService *ExchangeRateService) ListExchangeRates(ctx context.Context) (responses.ListExchangeRatesResponse, error) {
	config.Logger.Info("Service: Listing all exchange rates")

	result, err := exchangeRateService.exchangeRateRepository.ListExchangeRates(ctx)
	if err != nil {
		config.Logger.Error("Service: Failed to list exchange rates", "error", err.Error())
		return responses.ListExchangeRatesResponse{}, err
	}

	config.Logger.Info("Service: Successfully listed exchange rates", "total", result.Total)
	return result, nil
}
