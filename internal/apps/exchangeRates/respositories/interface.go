package exchangeRates

import (
	"context"

	requests "lemfi/simplebank/internal/apps/exchangeRates/requests"
	responses "lemfi/simplebank/internal/apps/exchangeRates/responses"
)

type ExchangeRateRepositoryInterface interface {
	ListExchangeRates(ctx context.Context) (responses.ListExchangeRatesResponse, error)
	GetExchangeRate(ctx context.Context, payload requests.GetExchangeRateRequest) (responses.GetExchangeRateResponse, error)
}
