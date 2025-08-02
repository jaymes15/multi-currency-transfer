package exchangeRates

import (
	"context"

	"lemfi/simplebank/config"
	requests "lemfi/simplebank/internal/apps/exchangeRates/requests"
	responses "lemfi/simplebank/internal/apps/exchangeRates/responses"
)

func (exchangeRateService *ExchangeRateService) GetExchangeRate(ctx context.Context, payload requests.GetExchangeRateRequest) (responses.GetExchangeRateResponse, error) {
	config.Logger.Info("Service: Getting exchange rate for currency pair",
		"from_currency", payload.FromCurrency,
		"to_currency", payload.ToCurrency,
		"amount", payload.Amount.String(),
	)

	result, err := exchangeRateService.exchangeRateRepository.GetExchangeRate(ctx, payload)
	if err != nil {
		config.Logger.Error("Service: Failed to get exchange rate", "error", err.Error())
		return responses.GetExchangeRateResponse{}, err
	}

	config.Logger.Info("Service: Successfully got exchange rate",
		"rate", result.ExchangeRate.Rate.String(),
		"can_transact", result.CanTransact,
	)

	return result, nil
}
